package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// NewHandlerContext gg
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

// UsersHandler gg
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method == http.MethodPost {
		// check for correct header
		ctype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "application/json") {
			// throw error
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		} else {
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
			anotherUser, err := dbUser.Insert(user)
			if err != nil {
				fmt.Errorf("Could not insert user to DB")
			}

			// ensure anotherUser contains the new database-assigned primary key value
			_, err = dbUser.GetByID(anotherUser.ID)
			if err != nil {
				fmt.Errorf("id does not contain the db primary key value")
			}

			userJSON, err := json.Marshal(anotherUser)
			if err != nil {
				fmt.Errorf("Could not marshal user")
			}
			// create a new session
			var sessionState SessionState
			sessionState.User = anotherUser
			sessionState.BeginTime = time.Now()

			_, err = sessions.BeginSession(ctx.Key, *ctx.Session, sessionState, w)
			if err != nil {
				fmt.Errorf("Could not begin session")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(userJSON)
		}
	}
}

// SessionsHandler ghgh
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect HTTP Method", http.StatusMethodNotAllowed)
	} else {
		// check for correct header
		ctype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "application/json") {
			// throw error
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		} else {
			fmt.Println(r.Body)
			var creds users.Credentials
			// jsonBody := r.Body

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

			userJSON, err := json.Marshal(user)
			if err != nil {
				fmt.Errorf("Could not marshal user")
			}
			// create a new session
			var sessionState SessionState
			sessionState.User = user
			sessionState.BeginTime = time.Now()

			_, err = sessions.BeginSession(ctx.Key, *ctx.Session, sessionState, w)
			if err != nil {
				fmt.Errorf("Could not begin session")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(userJSON)
		}
	}
}

// SpecificSessionsHandler ghgh
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
		// func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
		_, err := sessions.EndSession(r, ctx.Key, *ctx.Session)
		if err != nil {
			fmt.Errorf("Could not end session")
		}
		fmt.Println("signed out")
	}
}
