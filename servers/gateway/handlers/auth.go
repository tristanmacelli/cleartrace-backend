package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// UsersHandler does something
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	if !correctHeader(w, r) {
		return
	}
	var nu users.NewUser
	// make sure this json is valid
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&nu)
	if err != nil {
		panic(err)
	}
	user, err := nu.ToUser()
	if err != nil {
		http.Error(w, "Invalid user information", http.StatusUnprocessableEntity)
		return
	}

	// save user to database
	userStore := ctx.UserStore
	user, err = userStore.Insert(user)
	if err != nil {
		http.Error(w, "Could not save user", http.StatusInternalServerError)
		return
	}

	// ensure anotherUser contains the new database-assigned primary key value
	user, _ = userStore.GetByID(user.ID)
	// Unreachable (assuming we succeed to insert, there will be a user with the given ID)
	// if err != nil {
	// 	http.Error(w, "Could not find user", http.StatusInternalServerError)
	// 	return
	// }

	userJSON := encodeUser(user)
	ctx.beginSession(user, w)
	formatResponse(w, http.StatusCreated, userJSON)
}

// SpecificUserHandler does something
func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPatch {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	//Authentication process
	// Check the values in the authentication handler passed to the responsewriter
	// authValue := r.Header.Get("Authorization")
	// sessionID, err := sessions.GetSessionID(r, authValue)
	// if err != nil {
	// 	http.Error(w, "You are not authenticated", http.StatusUnauthorized)
	// 	return
	// }
	var userID []string = strings.Split(r.URL.String(), "users/")

	if r.Method == http.MethodGet {
		id, _ := strconv.ParseInt(userID[1], 10, 64)
		user, err := ctx.UserStore.GetByID(id)

		// We are checking for a nil value since our GetBy method will not return an
		// error if there were no matches (since this is not necessarily a failure)
		// We are also checking for an unpopulated user field because in the case that nobody
		// is found, the field will not be populated (in the case that someone is found the
		// user should always have that value)
		if user.FirstName == "" && err == nil {
			http.Error(w, "There is no user with the corresponding ID", http.StatusNotFound)
			return
		}
		userJSON := encodeUser(user)
		formatResponse(w, http.StatusOK, userJSON)
		// Patch
	} else {
		// if userID[1] != "me" && userID[1] != sessionID.String() {
		// 	http.Error(w, "You are unauthorized to perform this action", http.StatusForbidden)
		// 	return
		// }
		if !correctHeader(w, r) {
			return
		}
		var up users.Updates
		// make sure this json is valid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&up)
		if err != nil {
			panic(err)
		}

		user, err := ctx.UserStore.Update(1, &up)
		userJSON := encodeUser(user)
		formatResponse(w, http.StatusOK, userJSON)
	}
}

// SessionsHandler does something
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
	} else {
		if !correctHeader(w, r) {
			return
		}
		var creds users.Credentials

		// make sure this json is valid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&creds)
		if err != nil {
			panic(err)
		}

		user, err := ctx.UserStore.GetByEmail(creds.Email)
		// TODO: do something that would take about the same amount of time as authenticating
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		err = user.Authenticate(creds.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// log all successful user sign-in attempts
		ctx.UserStore.LogSuccessfulSignIns(user, r)

		userJSON := encodeUser(user)
		ctx.beginSession(user, w)
		formatResponse(w, http.StatusCreated, userJSON)
	}
}

// SpecificSessionsHandler does something
func (ctx *HandlerContext) SpecificSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != http.MethodDelete {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
	} else {
		mine := strings.HasSuffix(r.URL.String(), "mine")
		if !mine {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, ctx.Key, ctx.SessionStore)
		if err != nil {
			http.Error(w, "Could not find user", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("signed out"))
	}
}

func correctHeader(w http.ResponseWriter, r *http.Request) bool {
	// check for correct header
	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		// throw error
		http.Error(w, "must be JSON", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

func encodeUser(user *users.User) []byte {
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("Could not marshal user")
	}
	return userJSON
}

func (ctx *HandlerContext) beginSession(user *users.User, w http.ResponseWriter) {
	// create a new session
	var sessionState SessionState
	sessionState.User = user
	sessionState.BeginTime = time.Now()

	_, err := sessions.BeginSession(ctx.Key, ctx.SessionStore, sessionState, w)
	if err != nil {
		fmt.Errorf("Could not begin session")
	}
}

func formatResponse(w http.ResponseWriter, status int, userJSON []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(userJSON)
}
