package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server-side-mirror/servers/gateway/models/users"
	"server-side-mirror/servers/gateway/sessions"
	"strconv"
	"strings"
	"time"
)

// UsersHandler creates a new user, enters them into the users database and begins a session for them
func (ctx *HandlerContext) UsersHandler(response http.ResponseWriter, request *http.Request) {
	// check for POST
	if request.Method != http.MethodPost {
		http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	if !correctHeader(response, request) {
		return
	}
	var nu users.NewUser

	// make sure this json is valid
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&nu)
	if err != nil {
		http.Error(response, "Could not parse new user.", http.StatusUnprocessableEntity)
		return
	}

	user, err := nu.ToUser()
	if err != nil {
		http.Error(response, "Invalid user information", http.StatusNotAcceptable)
		return
	}

	// save user to database
	userStore := ctx.UserStore
	user, err = userStore.Insert(user)
	if err != nil {
		http.Error(response, "Could not save user", http.StatusInternalServerError)
		return
	}

	// Save user's ID to searchable index of user IDs
	userIndexes := ctx.UserIndexes
	userIndexes.Add(user.FirstName, user.ID)
	userIndexes.Add(user.LastName, user.ID)
	userIndexes.Add(user.UserName, user.ID)

	// ensure anotherUser contains the new database-assigned primary key value
	user, _ = userStore.GetByID(user.ID)
	userJSON := encodeUser(user)
	ctx.beginSession(user, response)
	formatResponse(response, http.StatusCreated, userJSON)
}

// SpecificUserHandler either returns all the user information or updates the users
// first and last names on a specific user
func (ctx *HandlerContext) SpecificUserHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		ctx.GetUserHandler(response, request)
		return
	}

	if request.Method == http.MethodPatch {
		ctx.UpdateUserHandler(response, request)
		return
	}

	http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
}

func (ctx *HandlerContext) GetUserHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}

	// Authentication process
	// Check the values in the authentication handler passed to the responsewriter

	sessionState := &SessionState{}
	_, err := sessions.GetState(request, ctx.Key, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(response, "You are not authenticated", http.StatusUnauthorized)
		return
	}
	var queryID []string = strings.Split(request.URL.String(), "users/")
	userID, _ := strconv.ParseInt(queryID[1], 10, 64)

	user, err := ctx.UserStore.GetByID(userID)
	if err != nil {
		http.Error(response, "Error getting user with the corresponding ID", http.StatusInternalServerError)
		return
	}
	// We are checking for a nil value since our GetBy method will not return an
	// error if there were no matches (since this is not necessarily a failure)
	// We are also checking for an unpopulated user field because in the case that nobody
	// is found, the field will not be populated (in the case that someone is found the
	// user should always have that value)
	if user.FirstName == "" && err == nil {
		http.Error(response, "There is no user with the corresponding ID", http.StatusNotFound)
		return
	}
	userJSON := encodeUser(user)
	formatResponse(response, http.StatusOK, userJSON)
}

func (ctx *HandlerContext) UpdateUserHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}

	// Authentication process
	// Check the values in the authentication handler passed to the responsewriter
	sessionState := &SessionState{}
	_, err := sessions.GetState(request, ctx.Key, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(response, "You are not authenticated", http.StatusUnauthorized)
		return
	}
	var user = sessionState.User
	var userID = user.ID
	var queryID []string = strings.Split(request.URL.String(), "users/")

	id, _ := strconv.ParseInt(queryID[1], 10, 64)
	if queryID[1] != "me" && id != userID {
		http.Error(response, "You are unauthorized to perform this action", http.StatusForbidden)
		return
	}
	if !correctHeader(response, request) {
		return
	}
	var up users.Updates
	// make sure this json is valid
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&up)
	if err != nil {
		panic(err)
	}
	// Save user's ID to searchable index of user IDs
	userIndexes := ctx.UserIndexes
	userIndexes.Remove(user.FirstName, userID)
	userIndexes.Remove(user.LastName, userID)
	userIndexes.Add(up.FirstName, userID)
	userIndexes.Add(up.LastName, userID)

	user, err = ctx.UserStore.Update(userID, &up)
	if err != nil {
		http.Error(response, "Failed to update user", http.StatusInternalServerError)
		return
	}
	userJSON := encodeUser(user)
	formatResponse(response, http.StatusOK, userJSON)
}

// SessionsHandler this logs a user in to our application (authentication)
func (ctx *HandlerContext) SessionsHandler(response http.ResponseWriter, request *http.Request) {
	// check for POST
	if request.Method != http.MethodPost {
		http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	if !correctHeader(response, request) {
		return
	}
	var creds users.Credentials
	// make sure this json is valid
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&creds)
	if err != nil {
		panic(err)
	}

	user, err := ctx.UserStore.GetByEmail(creds.Email)
	// TODO: do something that would take about the same amount of time as authenticating
	if err != nil {
		http.Error(response, "invalid credentials", http.StatusUnauthorized)
		return
	}

	err = user.Authenticate(creds.Password)
	if err != nil {
		http.Error(response, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// log all successful user sign-in attempts
	ctx.UserStore.LogSuccessfulSignIns(user, request)

	userJSON := encodeUser(user)
	ctx.beginSession(user, response)
	formatResponse(response, http.StatusCreated, userJSON)
}

// SpecificSessionsHandler logs a user out of our application
func (ctx *HandlerContext) SpecificSessionsHandler(response http.ResponseWriter, request *http.Request) {
	// check for POST
	if request.Method != http.MethodDelete {
		http.Error(response, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	mine := strings.HasSuffix(request.URL.String(), "mine")
	if !mine {
		http.Error(response, "access denied", http.StatusForbidden)
		return
	}
	_, err := sessions.EndSession(request, ctx.Key, ctx.SessionStore)
	if err != nil {
		http.Error(response, "Could not find user", http.StatusInternalServerError)
		return
	}
	// Send message to rabbitMQ indicating connection was removed
	message := mqMessage{}
	message.Type = "close-connection"
	// Error check message sending
	response.Write([]byte("signed out"))
}

func correctHeader(response http.ResponseWriter, request *http.Request) bool {
	// check for correct header
	ctype := request.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		// throw error
		http.Error(response, "must be JSON", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

func encodeUser(user *users.User) []byte {
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Could not marshal user: %s", err)
		panic("Could not encode user data")
	}
	return userJSON
}

func (ctx *HandlerContext) beginSession(user *users.User, response http.ResponseWriter) {
	// create a new session
	var sessionState SessionState
	sessionState.User = user
	sessionState.BeginTime = time.Now()
	_, err := sessions.BeginSession(ctx.Key, ctx.SessionStore, sessionState, response)
	if err != nil {
		fmt.Printf("Could not begin session")
		panic("Could not begin session")
	}
}

func formatResponse(response http.ResponseWriter, status int, userJSON []byte) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	response.Write(userJSON)
}
