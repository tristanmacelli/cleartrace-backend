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

// NewHandlerContext does something
func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	if user == nil {
		panic("No user")
	} else if session == nil {
		panic("No Session")
	} else if len(key) == 0 {
		panic("No User key")
	}
	return &HandlerContext{key, user, session}
}

// UsersHandler does something
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
		return
	}
	if correctHeader(w, r) {
		return
	}
	fmt.Println(r.Body)
	var nu users.NewUser
	// jsonBody := r.Body

	// make sure this json is valid
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&nu)
	if err != nil {
		panic(err)
	}
	// create a new user
	// nu.Email = jsonBody.Email
	// nu.Password = jsonBody.Password
	// nu.PasswordConf = jsonBody.PasswordConf
	// nu.UserName = jsonBody.UserName
	// nu.FirstName = jsonBody.FirstName
	// nu.LastNamw = jsonBody.LastName

	user, err := nu.ToUser()
	if err != nil {
		fmt.Errorf("Could not create a new user")
	}

	// save user to database
	dbUser := *ctx.User
	user, err = dbUser.Insert(user)
	if err != nil {
		fmt.Errorf("Could not insert user to DB")
	}

	// ensure anotherUser contains the new database-assigned primary key value
	_, err = dbUser.GetByID(user.ID)
	if err != nil {
		fmt.Errorf("id does not contain the db primary key value")
	}

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

	var userID []string = strings.Split(r.URL.String(), "users/")
	if r.Method == http.MethodGet {
		dbUser := *ctx.User
		id, _ := strconv.ParseInt(userID[1], 10, 64)
		user, err := dbUser.GetByID(id)

		if user.FirstName == "" && err == nil {
			http.Error(w, "There is no user with the corresponding ID", http.StatusNotFound)
			return
		}
		userJSON := encodeUser(user)
		formatResponse(w, http.StatusOK, userJSON)
		// Patch
	} else {
		// Conditionally convert the userID into an int to be able to compare with the corresponding authenticated user

		// 1 in the second case after the != is a mock of the currently authenticated user
		if userID[1] != "me" && userID[1] != 1 {
			http.Error(w, "You are unauthorized to perform this action", http.StatusForbidden)
			return
		}
		if !correctHeader(w, r) {
			return
		}
		var up users.Updates
		// jsonBody := r.Body

		// make sure this json is valid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&up)
		if err != nil {
			panic(err)
		}

		dbUser := *ctx.User
		user, err := dbUser.Update(1, &up)
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
		if correctHeader(w, r) {
			return
		}
		var creds users.Credentials

		// make sure this json is valid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&creds)
		if err != nil {
			panic(err)
		}

		dbUser := *ctx.User
		user, err := dbUser.GetByEmail(creds.Email)
		// TODO: do something that would take about the same amount of time as authenticating
		if err == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		err = user.Authenticate(creds.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

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
		_, err := sessions.EndSession(r, ctx.Key, *ctx.Session)
		if err != nil {
			fmt.Errorf("Could not end session")
		}
		// Is this the correct way to be Responding with the plain text message "signed out"
		fmt.Println("signed out")
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

	_, err := sessions.BeginSession(ctx.Key, *ctx.Session, sessionState, w)
	if err != nil {
		fmt.Errorf("Could not begin session")
	}
}

func formatResponse(w http.ResponseWriter, status int, userJSON []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(userJSON)
}
