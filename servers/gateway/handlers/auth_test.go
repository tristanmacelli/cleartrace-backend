package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestUserHandler(t *testing.T) {

    //email        string `json:"email"`
	// Password     string `json:"password"`
	// PasswordConf string `json:"passwordConf"`
	// UserName     string `json:"userName"`
	// FirstName    string `json:"firstName"`
	// LastName 

    jsonBody := []byte({   
        "Email":        "myexampleEmail@live.com",
        "Password":     "mypassword123",
        "PasswordConf": "mypassword123",
        "UserName":     "TMcGee123",
        "FirstName":    "Tester",
        "LastName":     "McGee"
    })

    req, err := http.NewRequest("POST", "v1/users/", bytes.NewBuffer(jsonBody))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(UserHandler)
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