package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server-side-mirror/servers/gateway/indexes"
	"server-side-mirror/servers/gateway/models/users"
	"server-side-mirror/servers/gateway/sessions"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var correctNewUser = map[string]string{
	"Email":        "myexampleEmail@live.com",
	"Password":     "mypassword123",
	"PasswordConf": "mypassword123",
	"UserName":     "TMcGee123",
	"FirstName":    "Tester",
	"LastName":     "McGee",
}

var incorrectNewUser = map[string]string{
	"Email":        "myexampleEmail@live.com",
	"Password":     "mypa",
	"PasswordConf": "mypassword123",
	"UserName":     "TMcGee123",
	"FirstName":    "Tester",
	"LastName":     "McGee",
}

var correctUserUpdates = map[string]string{
	"FirstName": "UpdatedTester",
	"LastName":  "UpdatedMcGee",
}

var incorrectUserUpdates = map[string]string{
	"Bad":      "Data",
	"UserName": "TMcGeeGee",
}

var correctCreds = map[string]string{
	"Email":    "myexampleEmail@live.com",
	"Password": "mypassword123",
}

var incorrectEmailCreds = map[string]string{
	"Email":    "myEmaillive.com",
	"Password": "mypassword123",
}

var incorrectPassCreds = map[string]string{
	"Email":    "myexampleEmail@live.com",
	"Password": "",
}

const headerAuthorization = "Authorization"
const schemeBearer = "Bearer "

func valueMapToUser(newUser map[string]string) *users.User {
	var nu users.NewUser
	nu.Email = newUser["Email"]
	nu.Password = newUser["Password"]
	nu.PasswordConf = newUser["PasswordConf"]
	nu.UserName = newUser["UserName"]
	nu.FirstName = newUser["FirstName"]
	nu.LastName = newUser["LastName"]
	user, _ := nu.ToUser()
	return user
}

// https://blog.questionable.services/article/testing-http-handlers-go/

// buildNewRequest creates a new request using the passed http method, path extras
// and value map as the json to be attached to the request
func buildNewRequest(
	t *testing.T,
	method string,
	contentType string,
	rawBodyData map[string]string,
	resourceIdentifier string,
	sessionID string,
) *http.Request {

	jsonBody, _ := json.Marshal(rawBodyData)
	path := "v1/users/" + resourceIdentifier
	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)
	authValue := schemeBearer + sessionID
	req.Header.Set(headerAuthorization, authValue)
	return req
}

// buildNewStores creates mock versions of the user and session stores for testing purposes
func buildStoresAndCtx(signingKey string) (sessions.Store, *HandlerContext) {
	userStore := &users.MockStore{}
	sessionStore := sessions.NewMemStore((time.Second * 20), (time.Second * 19))
	socketStore := NewNotify(map[int64]*websocket.Conn{}, &sync.Mutex{})
	indexedUsers := indexes.NewTrie(&sync.Mutex{})

	ctx := NewHandlerContext(signingKey, userStore, *indexedUsers, sessionStore, *socketStore)

	return sessionStore, ctx
}

// callUsersHandler calls the buildNewRequest and buildNewStores helper functions
// and calls the associated UsersHandler with mocked returns and errors for testing
func callUsersHandler(
	t *testing.T,
	method string,
	contentType string,
	rawUser map[string]string,
	err error,
) *httptest.ResponseRecorder {
	signingKey := "signing key"
	sessionID := "sessionid"
	_, ctx := buildStoresAndCtx(signingKey)
	rr := httptest.NewRecorder()

	users.SetErr(err)
	user := valueMapToUser(rawUser)
	users.SetInsertNextReturn(user)
	users.SetGetByIDNextReturn(user)

	req := buildNewRequest(t, method, contentType, rawUser, "", sessionID)
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// callSpecificUserHandler calls the buildNewRequest and buildNewStores helper functions
// and calls the associated SpecificUserHandler with mocked returns for testing
func callSpecificUserHandler(
	t *testing.T,
	method string,
	contentType string,
	rawUserUpdates map[string]string,
	resourceIdentifier string,
	user *users.User,
	err error,
	useExistingSessionID bool,
) *httptest.ResponseRecorder {
	signingKey := "signing key"
	sessionStore, ctx := buildStoresAndCtx(signingKey)
	rr := httptest.NewRecorder()

	users.SetErr(err)
	if method == "GET" {
		users.SetGetByIDNextReturn(user)
	}

	sessionState := SessionState{
		User:      user,
		BeginTime: time.Now(),
	}
	newSessionID, _ := sessions.BeginSession(signingKey, sessionStore, sessionState, rr)
	sessionID := ""
	if useExistingSessionID {
		sessionID = newSessionID.String()
	}

	req := buildNewRequest(t, method, contentType, rawUserUpdates, resourceIdentifier, sessionID)
	handler := http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// callSessionsHandler calls the buildNewRequest and buildNewStores helper functions
// and calls the associated SessionsHandler with mocked returns for testing
func callSessionsHandler(
	t *testing.T,
	method string,
	contentType string,
	credentials map[string]string,
	err error,
) *httptest.ResponseRecorder {
	signingKey := "signing key"
	sessionID := "sessionid"
	_, ctx := buildStoresAndCtx(signingKey)
	rr := httptest.NewRecorder()

	if err == nil {
		user := valueMapToUser(correctNewUser)
		users.SetGetByEmailNextReturn(user)
	}
	users.SetErr(err)

	req := buildNewRequest(t, method, contentType, credentials, "", sessionID)
	handler := http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// callSpecificSessionsHandler calls the buildNewRequest and buildNewStores helper functions
// and calls the associated SpecificSessionsHandler with mocked returns for testing
func callSpecificSessionsHandler(
	t *testing.T,
	method string,
	contentType string,
	resourceIdentifier string,
	useExistingSessionID bool,
) *httptest.ResponseRecorder {
	signingKey := "signing key"
	sessionStore, ctx := buildStoresAndCtx(signingKey)
	rr := httptest.NewRecorder()

	// Create a user and put it into session state
	sessionState := SessionState{
		User:      valueMapToUser(correctNewUser),
		BeginTime: time.Now(),
	}
	// Create a session for the user in session state
	newSessionID, _ := sessions.BeginSession(signingKey, sessionStore, sessionState, rr)
	sessionID := ""
	if useExistingSessionID {
		sessionID = newSessionID.String()
	}

	req := buildNewRequest(t, method, contentType, nil, resourceIdentifier, sessionID)
	handler := http.HandlerFunc(ctx.SpecificSessionsHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// TestUserHandler does something (5 total cases)
func TestUserHandler(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		method      string
		contentType string
		newUser     map[string]string
		err         error
		status      int
		expectation string
	}{
		{
			"POST method header set correctly",
			"Make sure you're using the correct http method",
			"POST",
			"application/json",
			correctNewUser,
			nil,
			http.StatusCreated,
			"expected a http.StatusCreated but the handler returned: %d",
		},
	}
	failcases := []struct {
		name        string
		hint        string
		method      string
		contentType string
		newUser     map[string]string
		err         error
		status      int
		expectation string
	}{
		{
			"GET method header not supported by this method",
			"Make sure you're using the correct http method",
			"GET",
			"application/json",
			correctNewUser,
			nil,
			http.StatusMethodNotAllowed,
			"expected a http.StatusMethodNotAllowed but the handler returned: %d",
		},
		{
			"Content-Type set incorrectly",
			"Make sure you're using the correct character encoding",
			"POST",
			"alication/json",
			correctNewUser,
			nil,
			http.StatusUnsupportedMediaType,
			"expected a http.StatusUnsupportedMediaType but the handler returned: %d",
		},
		{
			"Invalid new user passed",
			"Make sure you're using passing a valid new user",
			"POST",
			"application/json",
			incorrectNewUser,
			nil,
			http.StatusNotAcceptable,
			"expected a http.StatusNotAcceptable but the handler returned: %d",
		},
		{
			"Invalid database reference passed",
			"Make sure you're using passing a valid database reference",
			"POST",
			"application/json",
			correctNewUser,
			errors.New("Could not connect to db"),
			http.StatusInternalServerError,
			"expected a http.StatusInternalServerError but the handler returned: %d",
		},
	}
	for _, c := range cases {
		// SUCCESS CASE
		response := callUsersHandler(t, c.method, c.contentType, c.newUser, c.err)
		if status := response.Code; status != c.status {
			t.Errorf(c.expectation, status)
			t.Log(c.hint)
		}
	}
	// FAIL CASE
	for _, fc := range failcases {
		response := callUsersHandler(t, fc.method, fc.contentType, fc.newUser, fc.err)
		if status := response.Code; status != fc.status {
			t.Errorf(fc.expectation, status)
			t.Log(fc.hint)
		}
	}
}

// TestSpecificUserHandler does something
// Authorization dependent test cases (6) not operational of 13 total
func TestSpecificUserHandler(t *testing.T) {
	cases := []struct {
		name                 string
		hint                 string
		method               string
		contentType          string
		userUpdates          map[string]string
		resourceIdentifier   string
		user                 *users.User
		err                  error
		useExistingSessionID bool
		status               int
		expectation          string
	}{
		{
			"Success Case GET Method",
			"Ensure the method, content type, user id (resource id), user, and sessionID (boolean) are valid",
			"GET",
			"application/json",
			nil,
			"1234",
			valueMapToUser(correctNewUser),
			nil,
			true,
			http.StatusOK,
			"expected a http.StatusOK but the handler returned: %d",
		},
		{
			"Success Case PATCH Method",
			"Ensure the method, content type, user updates, user id (resource id), user, and sessionID (boolean) are valid",
			"PATCH",
			"application/json",
			correctUserUpdates,
			"me",
			valueMapToUser(correctNewUser),
			nil,
			true,
			http.StatusOK,
			"expected a http.StatusOK but the handler returned: %d",
		},
	}
	failcases := []struct {
		name                 string
		hint                 string
		method               string
		contentType          string
		userUpdates          map[string]string
		resourceIdentifier   string
		user                 *users.User
		err                  error
		useExistingSessionID bool
		status               int
		expectation          string
	}{
		{
			"Method header not allowed",
			"Must use either GET or PATCH http methods",
			"POST",
			"application/json",
			correctUserUpdates,
			"",
			valueMapToUser(correctNewUser),
			nil,
			true,
			http.StatusMethodNotAllowed,
			"expected a http.StatusMethodNotAllowed but the handler returned: %d",
		},
		{
			"Not user found",
			"There was no user found with the given id",
			"GET",
			"application/json",
			nil,
			"1234",
			&users.User{},
			errors.New("Not Authorized"),
			false,
			http.StatusUnauthorized,
			"expected a http.StatusUnauthorized but the handler returned: %d",
		},
		{
			"User ID Mismatch",
			"The passed userID does not match that of the current user",
			"PATCH",
			"application/json",
			correctUserUpdates,
			"1234",
			valueMapToUser(correctNewUser),
			nil,
			true,
			http.StatusForbidden,
			"expected a http.StatusForbidden but the handler returned: %d",
		},
		{
			"Content-Type set incorrectly",
			"Make sure you're using the correct character encoding",
			"PATCH",
			"alication/json",
			correctUserUpdates,
			"me",
			valueMapToUser(correctNewUser),
			nil,
			true,
			http.StatusUnsupportedMediaType,
			"expected a http.StatusUnsupportedMediaType but the handler returned: %d",
		},
		{
			"Invalid user updates format",
			"The passed user updates are not in a valid format",
			"PATCH",
			"application/json",
			incorrectUserUpdates,
			"me",
			valueMapToUser(correctNewUser),
			errors.New("Invalid user updates"),
			true,
			http.StatusInternalServerError,
			"522 expected a http.StatusInternalServerError but the handler returned: %d",
		},
	}
	for _, c := range cases {
		// SUCCESS CASE
		response := callSpecificUserHandler(
			t, c.method, c.contentType, c.userUpdates, c.resourceIdentifier, c.user, c.err, true,
		)
		if status := response.Code; status != c.status {
			t.Errorf(c.expectation, status)
			t.Log(c.hint)
		}
	}

	for _, fc := range failcases {
		// FAIL CASE
		response := callSpecificUserHandler(
			t, fc.method, fc.contentType, fc.userUpdates, fc.resourceIdentifier, fc.user, fc.err, fc.useExistingSessionID,
		)
		if status := response.Code; status != fc.status {
			t.Errorf(fc.expectation, status)
			t.Log(fc.hint)
		}
	}
}

// TestSessionsHandler does something
func TestSessionsHandler(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		method      string
		contentType string
		credentials map[string]string
		err         error
		status      int
		expectation string
	}{
		{
			"Success Case",
			"Ensure the method, content type, and credentials are valid",
			"POST",
			"application/json",
			correctCreds,
			nil,
			http.StatusCreated,
			"expected a http.StatusCreated but the handler returned: %d",
		},
	}
	failcases := []struct {
		name        string
		hint        string
		method      string
		contentType string
		credentials map[string]string
		err         error
		status      int
		expectation string
	}{
		{
			"PATCH method header not supported by this method",
			"Make sure you're using the correct http method",
			"PATCH",
			"application/json",
			correctCreds,
			nil,
			http.StatusMethodNotAllowed,
			"expected a http.StatusMethodNotAllowed but the handler returned: %d",
		},
		{
			"Content-Type set incorrectly",
			"Make sure you're using the correct character encoding",
			"POST",
			"alication/json",
			correctCreds,
			nil,
			http.StatusUnsupportedMediaType,
			"expected a http.StatuStatusUnsupportedMediaTypesMethodNotAllowed but the handler returned: %d",
		},
		{
			"Incorrect Credentials: Email",
			"Ensure the email is associated with a registered user",
			"POST",
			"application/json",
			incorrectEmailCreds,
			errors.New("Invalid Credentials, try again"),
			http.StatusUnauthorized,
			"expected a http.StatusUnauthorized but the handler returned: %d",
		},
		{
			"Incorrect Credentials: Password",
			"Ensure the password is valid",
			"POST",
			"application/json",
			incorrectPassCreds,
			errors.New("Invalid Credentials, try again"),
			http.StatusUnauthorized,
			"expected a http.StatusUnauthorized but the handler returned: %d",
		},
	}
	for _, c := range cases {
		// SUCCESS CASE
		response := callSessionsHandler(t, c.method, c.contentType, c.credentials, c.err)
		if status := response.Code; status != c.status {
			t.Errorf(c.expectation, status)
			t.Log(c.hint)
		}

	}
	for _, fc := range failcases {
		// FAIL CASE
		response := callSessionsHandler(t, fc.method, fc.contentType, fc.credentials, fc.err)
		if status := response.Code; status != fc.status {
			t.Errorf(fc.expectation, status)
			t.Log(fc.hint)
		}
	}
}

// TestSpecificSessionsHandler does something
func TestSpecificSessionsHandler(t *testing.T) {
	cases := []struct {
		name                 string
		hint                 string
		method               string
		contentType          string
		resourceIdentifier   string
		useExistingSessionID bool
		status               int
		expectation          string
	}{
		{
			"Success Case",
			"Ensure the method, content type, resource identifier, and sessionID (boolean) are valid",
			"DELETE",
			"application/json",
			"mine",
			true,
			http.StatusOK,
			"expected a http.StatusOK but the handler returned: %d",
		},
	}
	failureCases := []struct {
		name                 string
		hint                 string
		method               string
		contentType          string
		resourceIdentifier   string
		useExistingSessionID bool
		status               int
		expectation          string
	}{
		{
			"DELETE method header set incorrectly",
			"Make sure you're using the correct http method",
			"PATCH",
			"application/json",
			"mine",
			true,
			http.StatusMethodNotAllowed,
			"expected a http.StatusMethodNotAllowed but the handler returned: %d",
		},
		{
			"Resource Identifier set incorrectly",
			"Make sure you're using the correct resource identifier",
			"DELETE",
			"application/json",
			"",
			true,
			http.StatusForbidden,
			"expected a http.StatusForbidden but the handler returned: %d",
		},
		{
			"Pass a signing key that does not exist in sessions (EndSession Fails)",
			"Make sure you're passing an existing session id",
			"DELETE",
			"application/json",
			"mine",
			false,
			http.StatusInternalServerError,
			"expected a http.StatusInternalServerError but the handler returned: %d",
		},
	}
	for _, c := range cases {
		// SUCCESS CASE
		response := callSpecificSessionsHandler(t, c.method, c.contentType, c.resourceIdentifier, c.useExistingSessionID)
		if status := response.Code; status != c.status {
			t.Errorf(c.expectation, status)
			t.Log(c.hint)
		}
	}

	for _, fc := range failureCases {
		// FAIL CASE
		response := callSpecificSessionsHandler(t, fc.method, fc.contentType, fc.resourceIdentifier, fc.useExistingSessionID)
		if status := response.Code; status != fc.status {
			t.Errorf(fc.expectation, status)
			t.Log(fc.hint)
		}
	}
}
