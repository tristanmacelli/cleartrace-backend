package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

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
			dbUser := &ctx.User
			anotherUser, err := dbUser.Insert(user)
			if err != nil {
				fmt.Errorf("Could not insert user to DB")
			}

			// create a new session
			// BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
			sessionId, err := BeginSession(ctx.Key, ctx.User, ctx.Session, w)
			if err != nil {
				fmt.Errorf("Could not begin session")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

		}
	}
}
