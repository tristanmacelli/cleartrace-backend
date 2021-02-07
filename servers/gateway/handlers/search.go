package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const MaxReturnedUserIDs = 20

func (ctx *HandlerContext) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("X-User") == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
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
	// Format the response data
	userIDsJSON, err := json.Marshal(userIDs)
	if err != nil {
		fmt.Printf("Could not marshal indexes: %s", err)
	}
	formatResponse(w, http.StatusOK, userIDsJSON)
}
