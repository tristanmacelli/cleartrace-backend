package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server-side-mirror/servers/gateway/sessions"
)

const MaxReturnedUserIDs = 20

func (ctx *HandlerContext) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, "You are not authenticated", http.StatusUnauthorized)
		return
	}
	query, ok := r.URL.Query()["q"]
	if !ok || len(query[0]) < 1 {
		http.Error(w, "Must Pass Search Query", http.StatusBadRequest)
		return
	}
	// Find the user IDs
	userIndexes := ctx.UserIndexes
	userIDs := userIndexes.Find(query[0], MaxReturnedUserIDs)
	userStore := ctx.UserStore
	// var users []users.User

	// // How can we order these by FirstName ascending?
	// for _, id := range userIDs {
	// 	user, err := userStore.GetByID(id)
	// 	if err == nil {
	// 		users = append(users, *user)
	// 	}
	// }
	// Returns all user objects ordered by FirstName
	users, err := userStore.GetByIDs(userIDs, "FirstName")
	// Format the response data
	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("Could not marshal indexes: %s", err)
	}
	formatResponse(w, http.StatusOK, usersJSON)
}
