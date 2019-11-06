package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// var nu users.NewUser
// 	nu.Password = "mypassword123"
// 	nu.PasswordConf = "mypassword123"
// 	nu.UserName = "TMcGee123"
// 	nu.Email = "myexampleEmail@live.com"
// 	nu.FirstName = "Tester"
// 	nu.LastName = "McGee"
func TestUserHandler(t *testing.T) {

	//email        string `json:"email"`
	// Password     string `json:"password"`
	// PasswordConf string `json:"passwordConf"`
	// UserName     string `json:"userName"`
	// FirstName    string `json:"firstName"`
	// LastName

	valueMap := map[string]string{
		"Email":        "myexampleEmail@live.com",
		"Password":     "mypassword123",
		"PasswordConf": "mypassword123",
		"UserName":     "TMcGee123",
		"FirstName":    "Tester",
		"LastName":     "McGee",
	}
	jsonBody, _ := json.Marshal(valueMap)

	req, err := http.NewRequest("POST", "v1/users/", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	userStore := users.UserStore{}
	sessionStore := sessions.SessionStore{}

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore.Store, sessionStore.Store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// // ID        int64  `json:"id"`
	// // Email     string `json:"-"` //never JSON encoded/decoded
	// // PassHash  []byte `json:"-"` //never JSON encoded/decoded
	// // UserName  string `json:"userName"`
	// // FirstName string `json:"firstName"`
	// // LastName  string `json:"lastName"`
	// // PhotoURL

	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	//     t.Errorf("handler returned unexpected body: got %v want %v",
	//         rr.Body.String(), expected)
	// }

}

// func newSessionStore() (sessions.SessionID, error) {

// 	key := "test key"
// 	state := 100
// 	respRec := httptest.NewRecorder()
// 	sid, err := sessions.BeginSession(key, store, state, respRec)
// 	return sid, err
// }
