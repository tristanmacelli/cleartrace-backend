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
	user, err = dbUser.GetByID(user.ID)
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
	// Check the values in the authentication handler passed to the responsewriter
	authValue := r.Header.Get("Authorization")
	sessionID, err := sessions.GetSessionID(r, authValue)
	if err != nil {
		http.Error(w, "You are not authenticated", http.StatusUnauthorized)
		return
	}

	// Checking if the authenticated is in the redis db
	var sessionState SessionState
	session := *ctx.Session
	err = session.Get(sessionID, sessionState)
	if err != nil {
		http.Error(w, "You are not authenticated", http.StatusUnauthorized)
		return
	}

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
		if userID[1] != "me" && userID[1] != sessionID.String() {
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
		ctx.logSuccessfulSignIns(user, r)

		userJSON := encodeUser(user)
		ctx.beginSession(user, w)
		formatResponse(w, http.StatusCreated, userJSON)
	}
}

//logSuccessfulSignIns does something
func (ctx *HandlerContext) logSuccessfulSignIns(user *users.User, r *http.Request) {
	uid := user.ID
	timeOfSignIn := time.Now()
	clientIP := r.RemoteAddr
	ips := r.Header.Get("X-Forwarded-For")

	if len(ips) > 1 {
		clientIP = strings.Split(ips, ",")[0]
	} else if len(ips) == 1 {
		clientIP = ips
	}
	store := *ctx.User
	db := store.NewStore()
	tx, _ := db.DB.Begin()
	insq := "INSERT INTO userSignIn(userID, signinDT, ip) VALUES (?,?,?)"
	_, err := tx.Exec(insq, uid, timeOfSignIn, clientIP)

	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		// Close the reserved connection upon failure
		tx.Rollback()
		return
	}
	// Close the reserved connection upon success
	tx.Commit()
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
