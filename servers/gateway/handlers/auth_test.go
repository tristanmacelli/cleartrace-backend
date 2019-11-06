package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	"Email":        "myexampleEmail@live.com",
	"Password":     "mypassword123",
	"PasswordConf": "mypassword123",
	"UserName":     "",
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

func buildRequest(t *testing.T, method string, contentType string, valueMap map[string]string, pathExtras string) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(valueMap)

	path := "v1/users/" + pathExtras
	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	authValue := schemeBearer + sessionID
	req.Header.Set(headerAuthorization, authValue)

	// db, moc, err := sqlmock.New()
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/demo", os.Getenv("MYSQL_ROOT_PASSWORD"))
	userStore := users.NewMysqlStore(dsn)
	// Add fields to this after running docker container to run tests
	sessionStore := sessions.RedisStore{}

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, &sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

func TestUserHandler(t *testing.T) {

	rr := buildRequest(t, "POST", "", correctNewUser, "")
	// Success Case
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildRequest(t, "GET", "", correctNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code")
	}

	rr = buildRequest(t, "POST", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnsupportedMediaType {
		t.Errorf(
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code")
	}

	rr = buildRequest(t, "POST", "alication/json", correctNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	}

	rr = buildRequest(t, "POST", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnprocessableEntity {
		t.Errorf(
			"we did not expect a http.StatusUnprocessableEntity but the handler returned this status code")
	}

	rr = buildRequest(t, "POST", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf(
			"we expected an http.StatusUnprocessableEntity but the handler returned wrong status code")
	}
}

// TestSpecificUserHandler does something
func TestSpecificUserHandler(t *testing.T) {
	rr := buildRequest(t, "GET", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}
	rr = buildRequest(t, "PATCH", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildRequest(t, "POST", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestSessionsHandler(t *testing.T) {
	rr := buildRequest(t, "POST", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildRequest(t, "PATCH", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code")
	}

	rr = buildRequest(t, "POST", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnsupportedMediaType {
		t.Errorf(
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code")
	}

	rr = buildRequest(t, "POST", "alication/json", correctNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	}

	rr = buildRequest(t, "POST", "application/json", correctCreds, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf(
			"we did not expect a http.StatusUnauthorized but the handler returned this status code")
	}

	rr = buildRequest(t, "POST", "application/json", incorrectEmailCreds, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf(
			"we expected an http.StatusUnauthorized but the handler returned wrong status code")
	}

	rr = buildRequest(t, "POST", "application/json", incorrectPassCreds, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf(
			"we expected an http.StatusUnauthorized but the handler returned wrong status code")
	}
}

// TestSpecificSessionsHandler does something
func TestSpecificSessionsHandler(t *testing.T) {
	rr := buildRequest(t, "DELETE", "application/json", correctNewUser, "")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf(
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code")
	}

	rr = buildRequest(t, "PATCH", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected a http.StatusMethodNotAllowed but the handler did not return this status code")
	}

	rr = buildRequest(t, "DELETE", "application/json", correctNewUser, "mine")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusForbidden {
		t.Errorf(
			"we did not expect a http.StatusForbidden but the handler returned this status code")
	}

	rr = buildRequest(t, "DELETE", "alication/json", incorrectNewUser, "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf(
			"we expected a http.StatusForbidden but the handler did not return this status code")
	}
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
