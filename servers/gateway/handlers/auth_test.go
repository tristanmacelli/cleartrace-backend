package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewUser Formats
// var nu users.NewUser
// 	nu.Password = "mypassword123"
// 	nu.PasswordConf = "mypassword123"
// 	nu.UserName = "TMcGee123"
// 	nu.Email = "myexampleEmail@live.com"
// 	nu.FirstName = "Tester"
// 	nu.LastName = "McGee"

//email        string `json:"email"`
// Password     string `json:"password"`
// PasswordConf string `json:"passwordConf"`
// UserName     string `json:"userName"`
// FirstName    string `json:"firstName"`
// LastName

// User format
// ID        int64  `json:"id"`
// Email     string `json:"-"` //never JSON encoded/decoded
// PassHash  []byte `json:"-"` //never JSON encoded/decoded
// UserName  string `json:"userName"`
// FirstName string `json:"firstName"`
// LastName  string `json:"lastName"`
// PhotoURL

var correctNewUser = map[string]string{
	"Email":        "myexampleEmail@live.com",
	"Password":     "mypassword123",
	"PasswordConf": "mypassword123",
	"UserName":     "TMcGee123",
	"FirstName":    "Tester",
	"LastName":     "McGee",
}

var incorrectNewUser = map[string]string{
	"Email":        "#@%^%#$@#$@#.com",
	"Password":     "mypassword123",
	"PasswordConf": "mypassword123",
	"UserName":     "TMcGee123",
	"FirstName":    "Tester",
	"LastName":     "McGee",
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
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

// TODO: Make sure all cases pass for TestUserHandler
// before moving to TestSpecificUserHandler
const sessionID = "1234"

func buildNewRequest(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string, sessionID string) *http.Request {

	jsonBody, _ := json.Marshal(valueMap)

	path := "v1/users/" + pathExtras
	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	authValue := schemeBearer + sessionID
	req.Header.Set(headerAuthorization, authValue)
	return req
}

func buildNewStores() (users.Store, sessions.Store) {
	// If we are going to use a concrete user store to test, we need to use the values we
	// generate from the docker run command to get database-name & MYSQL_ROOT_PASSWORD
	// dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/database-name", os.Getenv("MYSQL_ROOT_PASSWORD"))
	ustore := users.MockStore{}
	var userStore users.Store
	userStore = &ustore
	// Add fields to this after running docker container to run tests
	sStore := sessions.MemStore{}
	var sessionStore sessions.Store
	sessionStore = &sStore
	return userStore, sessionStore
}

func buildCtxUser(t *testing.T, method string, contentType string,
	valueMap map[string]string, err bool) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, "", "1234")
	userStore, sessionStore := buildNewStores()

	if err {
		users.SetErr(errors.New("Could not insert"))
	} else {
		users.SetErr(nil)
	}
	var nu users.NewUser
	nu.Email = valueMap["Email"]
	nu.Password = valueMap["Password"]
	nu.PasswordConf = valueMap["PasswordConf"]
	nu.UserName = valueMap["UserName"]
	nu.FirstName = valueMap["FirstName"]
	nu.LastName = valueMap["LastName"]
	user, _ := nu.ToUser()
	users.SetInsertnextReturn(user)

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

func buildCtxSpecificUser(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string, sessionID string) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, pathExtras, sessionID)
	userStore, sessionStore := buildNewStores()

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

func buildCtxSession(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, pathExtras, "1234")
	userStore, sessionStore := buildNewStores()

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

func buildCtxSpecificSession(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, pathExtras, "1234")
	userStore, sessionStore := buildNewStores()

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SpecificSessionsHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// TestUserHandler does something
// TODO: Check if we need getbyid cases
// NOTE: 1 case outside of stores fails
func TestUserHandler(t *testing.T) {

	rr := buildCtxUser(t, "POST", "", correctNewUser, false)
	// Success Case
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildCtxUser(t, "GET", "", correctNewUser, false)
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code")
	}

	rr = buildCtxUser(t, "POST", "application/json", correctNewUser, false)
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnsupportedMediaType {
		t.Errorf(
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code")
	}

	rr = buildCtxUser(t, "POST", "alication/json", correctNewUser, false)
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	}

	rr = buildCtxUser(t, "POST", "application/json", correctNewUser, false)
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnprocessableEntity {
		t.Errorf(
			"we did not expect a http.StatusUnprocessableEntity but the handler returned this status code")
	}

	// rr = buildCtxUser(t, "POST", "application/json", incorrectNewUser, false)
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusUnprocessableEntity {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code: %v",
	// 		status)
	// }

	// Need test cases for INSERT
	rr = buildCtxUser(t, "POST", "application/json", correctNewUser, false)
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusInternalServerError {
		t.Errorf(
			"we did not expect a http.StatusInternalServerError but the handler returned this status code")
	}

	// // Pass incorrect dsn/invalid store reference
	rr = buildCtxUser(t, "POST", "application/json", correctNewUser, true)
	// FAIL CASE
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf(
			"we expected an http.StatusInternalServerError but the handler returned wrong status code %v", status)
	}

	// Test cases for GetByID
	// rr = buildCtxUser(t, "POST", "application/json", correctNewUser)
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusInternalServerError {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusInternalServerError but the handler returned this status code")
	// }

	// rr = buildCtxUser(t, "POST", "application/json", correctNewUser)
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusInternalServerError {
	// 	t.Errorf(
	// 		"we expected an http.StatusInternalServerError but the handler returned wrong status code")
	// }
}

// TestSpecificUserHandler does something
// All test cases written
func TestSpecificUserHandler(t *testing.T) {
	rr := buildCtxSpecificUser(t, "GET", "application/json", correctNewUser, "", "1234")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}
	rr = buildCtxSpecificUser(t, "PATCH", "application/json", correctNewUser, "", "1234")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildCtxSpecificUser(t, "POST", "alication/json", incorrectNewUser, "", "1234")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	// Need test cases for GetSessionID
	// passing sessionid in ctx that does exist in our sessions
	rr = buildCtxSpecificUser(t, "GET", "application/json", correctNewUser, "1234", "1234")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf(
			"we did not expect a http.StatusNotFound but the handler returned this status code")
	}

	// passing sessionid in ctx that does not exist in our sessions
	rr = buildCtxSpecificUser(t, "GET", "alication/json", correctNewUser, "123", "1234")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf(
			"we expected an http.StatusNotFound but the handler returned wrong status code")
	}

	// In If branch
	// Need test cases for GetByID when using GET method
	// passing sessionid in path does exist in our sessions
	// rr = buildCtxSpecificUser(t, "GET", "application/json", correctNewUser, "1234")
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusNotFound {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusNotFound but the handler returned this status code")
	// }

	// // passing sessionid in path that does not exist in our sessions
	// rr = buildCtxSpecificUser(t, "GET", "alication/json", correctNewUser, "123")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusNotFound {
	// 	t.Errorf(
	// 		"we expected an http.StatusNotFound but the handler returned wrong status code")
	// }

	// In else branch
	// Need test cases for authenticated OR matching sessionID
	// TODO: Refactor build request to accept an sessionID to handle this testing
	// rr = buildCtxSpecificUser(t, "PATCH", "application/json", correctNewUser, "1234")
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusForbidden {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusForbidden but the handler returned this status code")
	// }

	// rr = buildCtxSpecificUser(t, "PATCH", "application/json", correctNewUser, "me")
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusForbidden {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusForbidden but the handler returned this status code")
	// }

	// // User is authorized, but not allowed to access user id 123
	// rr = buildCtxSpecificUser(t, "PATCH", "alication/json", correctNewUser, "123")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusUnsupportedMediaType {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	// }

	// // malformed path
	// rr = buildCtxSpecificUser(t, "PATCH", "alication/json", correctNewUser, "m")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusUnsupportedMediaType {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	// }

	// // Checking for correct headers
	// rr = buildCtxSpecificUser(t, "PATCH", "application/json", correctNewUser, "")
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusUnsupportedMediaType {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code")
	// }

	// rr = buildCtxSpecificUser(t, "PATCH", "alication/json", correctNewUser, "")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusUnsupportedMediaType {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	// }

}

// TestSessionsHandler does something
// All test cases written
func TestSessionsHandler(t *testing.T) {
	rr := buildCtxSession(t, "POST", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildCtxSession(t, "PATCH", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code")
	}

	rr = buildCtxSession(t, "POST", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnsupportedMediaType {
		t.Errorf(
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code")
	}

	rr = buildCtxSession(t, "POST", "alication/json", correctNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	}

	// rr = buildCtxSession(t, "POST", "application/json", correctCreds, "")
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusUnauthorized {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusUnauthorized but the handler returned this status code")
	// }

	// rr = buildCtxSession(t, "POST", "application/json", incorrectEmailCreds, "")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusUnauthorized {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnauthorized but the handler returned wrong status code")
	// }

	// rr = buildCtxSession(t, "POST", "application/json", incorrectPassCreds, "")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusUnauthorized {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnauthorized but the handler returned wrong status code")
	// }
}

// TestSpecificSessionsHandler does something
// TODO: Write 2 test cases
func TestSpecificSessionsHandler(t *testing.T) {
	rr := buildCtxSpecificSession(t, "DELETE", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildCtxSpecificSession(t, "PATCH", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected a http.StatusMethodNotAllowed but the handler did not return this status code")
	}

	rr = buildCtxSpecificSession(t, "DELETE", "application/json", correctNewUser, "mine")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusForbidden {
		t.Errorf(
			"we did not expect a http.StatusForbidden but the handler returned this status code")
	}

	rr = buildCtxSpecificSession(t, "DELETE", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf(
			"we expected a http.StatusForbidden but the handler did not return this status code")
	}

	// Need test cases for EndSession
	// Pass a signing key that exists in sessions
	// rr = buildCtxSpecificSession(t, "DELETE", "application/json", correctNewUser, "mine")
	// // SUCCESS CASE
	// if status := rr.Code; status == http.StatusInternalServerError {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusInternalServerError but the handler returned this status code")
	// }

	// // Pass a signing key that does not exist in sessions
	// rr = buildCtxSpecificSession(t, "DELETE", "application/json", incorrectNewUser, "mine")
	// // FAIL CASE
	// if status := rr.Code; status != http.StatusInternalServerError {
	// 	t.Errorf(
	// 		"we expected a http.StatusInternalServerError but the handler did not return this status code")
	// }
}

// Random Comments

// expected := `{"alive": true}`
// if rr.Body.String() != expected {
//     t.Errorf("handler returned unexpected body: got %v want %v",
//         rr.Body.String(), expected)
// }

// userStore := users.UserStore{}
// sessionStore := sessions.SessionStore{}

// func newSessionStore() (sessions.SessionID, error) {

// 	key := "test key"
// 	state := 100
// 	respRec := httptest.NewRecorder()
// 	sid, err := sessions.BeginSession(key, store, state, respRec)
// 	return sid, err
// }
