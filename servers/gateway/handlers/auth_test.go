package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

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

// ID        int64  `json:"id"`
// Email     string `json:"-"` //never JSON encoded/decoded
// PassHash  []byte `json:"-"` //never JSON encoded/decoded
// UserName  string `json:"userName"`
// FirstName string `json:"firstName"`
// LastName  string `json:"lastName"`
// PhotoURL

var valueMap = map[string]string{
	"Email":        "myexampleEmail@live.com",
	"Password":     "mypassword123",
	"PasswordConf": "mypassword123",
	"UserName":     "TMcGee123",
	"FirstName":    "Tester",
	"LastName":     "McGee",
}

func buildRequest(t *testing.T, method string, contentType string) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(valueMap)

	req, err := http.NewRequest(method, "v1/users/", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	db, _, err := sqlmock.New()
	userStore := users.MysqlStore{}
	userStore.DB = db
	sessionStore := sessions.RedisStore{}

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", &userStore, &sessionStore)
	rr := httptest.NewRecorder()
	// ctx.SessionsHandler(rr ,req)
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

func TestUserHandler(t *testing.T) {
	rr := buildRequest(t, "POST", "")

	// Success Case
	if status := rr.Code; status == http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	rr = buildRequest(t, "GET", "")
	// FAIL CASE
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	rr = buildRequest(t, "POST", "application/json")
	// SUCCESS CASE
	if status := rr.Code; status == http.StatusUnsupportedMediaType {
		t.Errorf(
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code")
	}

	rr = buildRequest(t, "POST", "alication/json")
	// FAIL CASE
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code")
	}
}

// func TestUserHandler(t *testing.T) {

// 	// FAIL CASE
// 	if status := rr.Code; status != http.StatusMethodNotAllowed {
// 		t.Errorf(
// 			"we expected an http.StatusMethodNotAllowed but the handler returned wrong status code: got %v want %v",
// 			status, http.StatusMethodNotAllowed)
// 	}
// }

// expected := `{"alive": true}`
// if rr.Body.String() != expected {
//     t.Errorf("handler returned unexpected body: got %v want %v",
//         rr.Body.String(), expected)
// }

// userStore := users.UserStore{}
// sessionStore := sessions.SessionStore{}

// // func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
// ctx := NewHandlerContext("anything", userStore.Store, sessionStore.Store)

// func newSessionStore() (sessions.SessionID, error) {

// 	key := "test key"
// 	state := 100
// 	respRec := httptest.NewRecorder()
// 	sid, err := sessions.BeginSession(key, store, state, respRec)
// 	return sid, err
// }
