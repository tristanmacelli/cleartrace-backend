package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server-side-mirror/servers/gateway/sessions"
)

const MaxReturnedUserIDs = 20

func (ctx *HandlerContext) SearchHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet && request.Method != http.MethodPost {
		http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	sessionState := &SessionState{}
	_, err := sessions.GetState(request, ctx.Key, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(response, "You are not authenticated", http.StatusUnauthorized)
		return
	}

	var userIDs []int64
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&userIDs)

	if request.Method == http.MethodPost && err != nil {
		http.Error(response, "Failed to unmarshall userID data", http.StatusUnprocessableEntity)
		return
	}

	query, ok := request.URL.Query()["q"]
	if request.Method == http.MethodGet {
		fmt.Println("Search Query:", query[0])
		if !ok || len(query[0]) < 1 {
			http.Error(response, "Must Pass Search Query", http.StatusBadRequest)
			return
		}
		// Find the user IDs
		userIndexes := ctx.UserIndexes
		userIDs, _ = userIndexes.Find(query[0], MaxReturnedUserIDs)
	}
	userStore := ctx.UserStore

	// Returns all user objects ordered by FirstName & then Lastname (if necessary)
	users, err := userStore.GetByIDs(userIDs, []string{"FirstName", "Lastname"})

	if err != nil {
		http.Error(response, "Failed to query users from database", http.StatusInternalServerError)
		return
	}
	// Format the response data
	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("Could not marshal indexes: %s", err)
	}
	formatResponse(response, http.StatusOK, usersJSON)
}
