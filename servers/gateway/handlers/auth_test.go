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
	"time"
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

// buildNewRequest creates a new request using the passed http method, path extras
// and value map as the json to be attached to the request
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

// buildNewStores creates mock versions of the user and session stores for testing purposes
func buildNewStores() (users.Store, sessions.Store) {
	ustore := users.MockStore{}
	var userStore users.Store
	userStore = &ustore
	sStore := sessions.NewMemStore((time.Second * 20), (time.Second * 19))
	var sessionStore sessions.Store
	sessionStore = sStore
	return userStore, sessionStore
}

// buildCtxUser calls the buildNewRequest and buildNewStores helper functions
// and calls the associated UsersHandler with mocked returns and errors for testing
func buildCtxUser(t *testing.T, method string, contentType string,
	valueMap map[string]string, expectedErr bool) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, "", "1234")
	userStore, sessionStore := buildNewStores()

	if expectedErr {
		users.SetErr(errors.New("Could not connect to db"))
	} else {
		users.SetErr(nil)
	}
	user := valueMapToUser(valueMap)
	users.SetInsertNextReturn(user)
	users.SetGetByIDNextReturn(user)

	ctx := NewHandlerContext("This should be a valid key", userStore, sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// buildCtxSpecificUser calls the buildNewRequest and buildNewStores helper functions
// and calls the associated SpecificUserHandler with mocked returns for testing
func buildCtxSpecificUser(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string, sessionID string,
	foundUser bool, expectedErr bool) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, pathExtras, sessionID)
	userStore, sessionStore := buildNewStores()

	// https://blog.questionable.services/article/testing-http-handlers-go/

	if expectedErr {
		users.SetErr(errors.New("Attn: No User"))
	} else {
		users.SetErr(nil)
	}
	user := valueMapToUser(valueMap)
	if method == "GET" && foundUser {
		users.SetGetByIDNextReturn(user)
	} else if method == "GET" {
		users.SetGetByIDNextReturn(&users.User{})
	}
	var sessionState SessionState
	sessionState.User = user
	sessionState.BeginTime = time.Now()

	sid, _ := sessions.NewSessionID(sessionID)
	sessionStore.Save(sid, sessionState)

	ctx := NewHandlerContext("1234", userStore, sessionStore)
	rr := httptest.NewRecorder()
	sessions.BeginSession("1234", sessionStore, sessionState, rr)
	handler := http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// buildCtxSession calls the buildNewRequest and buildNewStores helper functions
// and calls the associated SessionsHandler with mocked returns for testing
func buildCtxSession(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string, expectedErr bool) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, pathExtras, "1234")
	userStore, sessionStore := buildNewStores()

	if expectedErr {
		users.SetErr(errors.New("Invalid Credentials, try again"))
	} else {
		user := valueMapToUser(correctNewUser)
		users.SetGetByEmailNextReturn(user)
		users.SetErr(nil)
	}

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, sessionStore)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// buildCtxSpecificSession calls the buildNewRequest and buildNewStores helper functions
// and calls the associated SpecificSessionsHandler with mocked returns for testing
func buildCtxSpecificSession(t *testing.T, method string, contentType string,
	valueMap map[string]string, pathExtras string) *httptest.ResponseRecorder {

	req := buildNewRequest(t, method, contentType, valueMap, pathExtras, "1234")
	userStore, sessionStore := buildNewStores()

	user := valueMapToUser(valueMap)
	var sessionState SessionState
	sessionState.User = user
	sessionState.BeginTime = time.Now()

	// func NewHandlerContext(key string, user *users.Store, session *sessions.Store) *HandlerContext {
	ctx := NewHandlerContext("anything", userStore, sessionStore)
	rr := httptest.NewRecorder()
	_, _ = sessions.BeginSession("anything", sessionStore, sessionState, rr)
	handler := http.HandlerFunc(ctx.SpecificSessionsHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// TestUserHandler does something
// TODO: Check if we need getbyid cases
// All tests pass! (8 total)
func TestUserHandler(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		method      string
		encoding    string
		nu          map[string]string
		status      int
		expected    string
		expectedErr bool
	}{
		{
			"POST method header set correctly",
			"Make sure you're using the correct http method",
			"POST",
			"",
			correctNewUser,
			http.StatusMethodNotAllowed,
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code",
			false,
		},
		{
			"Content-Type set correctly",
			"Make sure you're using the correct character encoding",
			"POST",
			"application/json",
			correctNewUser,
			http.StatusUnsupportedMediaType,
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code",
			false,
		},
		{
			"Correct new user passed",
			"Make sure you're passing a new user with valid data",
			"POST",
			"application/json",
			correctNewUser,
			http.StatusUnprocessableEntity,
			"we did not expect a http.StatusUnprocessableEntity but the handler returned this status code",
			false,
		},
		{
			"User successfully saved to the database",
			"Make sure you're passing a valid database reference",
			"POST",
			"application/json",
			correctNewUser,
			http.StatusInternalServerError,
			"we did not expect a http.StatusInternalServerError but the handler returned this status code",
			false,
		},
	}
	failcases := []struct {
		name        string
		hint        string
		method      string
		encoding    string
		nu          map[string]string
		status      int
		expected    string
		expectedErr bool
	}{
		{
			"GET method header not supported by this method",
			"Make sure you're using the correct http method",
			"GET",
			"",
			correctNewUser,
			http.StatusMethodNotAllowed,
			"we expected a http.StatusMethodNotAllowed but the handler did not return this status code",
			false,
		},
		{
			"Content-Type set incorrectly",
			"Make sure you're using the correct character encoding",
			"POST",
			"alication/json",
			correctNewUser,
			http.StatusUnsupportedMediaType,
			"we expected a http.StatusUnsupportedMediaType but the handler did not return this status code",
			false,
		},
		{
			"Invalid new user passed",
			"Make sure you're using passing a valid new user",
			"POST",
			"application/json",
			incorrectNewUser,
			http.StatusUnprocessableEntity,
			"we expected a http.StatusUnprocessableEntity but the handler did not return this status code",
			false,
		},
		{
			"Invalid database reference passed",
			"Make sure you're using passing a valid database reference",
			"POST",
			"application/json",
			correctNewUser,
			http.StatusInternalServerError,
			"we expected a http.StatusInternalServerError but the handler did not return this status code",
			true,
		},
	}
	for i, c := range cases {
		// SUCCESS CASE
		rr := buildCtxUser(t, c.method, c.encoding, c.nu, c.expectedErr)
		if status := rr.Code; status == c.status {
			t.Errorf(c.expected)
		}
		// FAIL CASE
		fc := failcases[i]
		rr = buildCtxUser(t, fc.method, fc.encoding, fc.nu, fc.expectedErr)
		if status := rr.Code; status != fc.status {
			t.Errorf(fc.expected)
		}
	}
}

// TestSpecificUserHandler does something
// Authorization dependent test cases (6) not operational of 13 total
func TestSpecificUserHandler(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		method      string
		encoding    string
		nu          map[string]string
		pathExtras  string
		sessionID   string
		foundUser   bool
		expectedErr bool
		status      int
		expected    string
	}{
		{
			"GET method header set correctly",
			"Make sure you're using the correct http method",
			"GET",
			"application/json",
			correctNewUser,
			"",
			"1234",
			true,
			false,
			http.StatusMethodNotAllowed,
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code",
		},
		{
			"Method header set correctly",
			"Make sure you're using the correct http method",
			"PATCH",
			"application/json",
			correctNewUser,
			"",
			"1234",
			true,
			false,
			http.StatusMethodNotAllowed,
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code",
		},
		{
			"User data was found (in If statement)",
			"Make sure you have signed up before trying to log-in",
			"GET",
			"application/json",
			correctNewUser,
			"1234",
			"1234",
			true,
			false,
			http.StatusNotFound,
			"we did not expect a http.StatusNotFound but the handler returned this status code",
		},
		{
			"Correctly updating user information (first & last names)",
			"Make sure you passing valid update format",
			"PATCH",
			"application/json",
			correctNewUser,
			"1234",
			"1234",
			true,
			false,
			http.StatusUnsupportedMediaType,
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code",
		},
	}
	failcases := []struct {
		name        string
		hint        string
		method      string
		encoding    string
		nu          map[string]string
		pathExtras  string
		sessionID   string
		foundUser   bool
		expectedErr bool
		status      int
		expected    string
	}{
		{
			"Method header not allowed",
			"Must use either GET or PATCH http methods",
			"POST",
			"alication/json",
			incorrectNewUser,
			"",
			"1234",
			true,
			false,
			http.StatusMethodNotAllowed,
			"we expected a http.StatusMethodNotAllowed but the handler did not return this status code",
		},
		{
			"No user found",
			"There was no user found with the given id",
			"GET",
			"alication/json",
			incorrectNewUser,
			"123",
			"1234",
			false,
			false,
			http.StatusNotFound,
			"we expected a http.StatusNotFound but the handler did not return this status code",
		},
		{
			"Invalid user updates format",
			"The passed user updates are not in a valid format",
			"PATCH",
			"alication/json",
			correctNewUser,
			"me",
			"me",
			true,
			true,
			http.StatusUnsupportedMediaType,
			"we expected a http.StatusUnsupportedMediaType but the handler did not return this status code",
		},
	}
	for i, fc := range failcases {
		// SUCCESS CASE
		c := cases[i]
		rr := buildCtxSpecificUser(t, c.method, c.encoding, c.nu, c.pathExtras, c.sessionID,
			c.foundUser, c.expectedErr)
		if status := rr.Code; status == c.status {
			t.Errorf(c.expected)
		}
		// FAIL CASE
		rr = buildCtxSpecificUser(t, fc.method, fc.encoding, fc.nu, fc.pathExtras, c.sessionID,
			c.foundUser, c.expectedErr)
		if status := rr.Code; status != fc.status {
			t.Errorf(fc.expected)
		}
	}
	c := cases[len(cases)-1]
	rr := buildCtxSpecificUser(t, c.method, c.encoding, c.nu, c.pathExtras, c.sessionID,
		c.foundUser, c.expectedErr)
	if status := rr.Code; status == c.status {
		t.Errorf(c.expected)
	}

	// Method header checks

	// Test cases for GetSessionID
	// THESE CURRENTLY DO NOT WORK BECAUSE WE WERE SURE HOW TO SAVE TO THE SESSIONID TO
	// THE SESSION STORE FOR TESTING PURPOSES.
	// passing sessionid in ctx that does exist in our sessions
	// // SUCCESS CASE
	// rr = buildCtxSpecificUser(t, "GET", "application/json", correctNewUser, "1234", "1234", true, false)
	// if status := rr.Code; status == http.StatusUnauthorized {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusUnauthorized but the handler returned this status code: %v",
	// 		status)
	// }

	// // passing sessionid in ctx that does not exist in our sessions
	// // FAIL CASE
	// rr = buildCtxSpecificUser(t, "GET", "application/json", correctNewUser, "123", "1234", true, false)
	// if status := rr.Code; status != http.StatusUnauthorized {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnauthorized but the handler returned wrong status code")
	// }

	// In If branch
	// Test cases for GetByID when using GET method

	// In else branch
	// Need test cases for authenticated OR matching sessionID
	// TODO: Refactor build request to accept an sessionID to handle this testing
	// // SUCCESS CASE
	// rr = buildCtxSpecificUser(t, "PATCH", "application/json", correctNewUser, "1234", "1234", true, false)
	// if status := rr.Code; status == http.StatusForbidden {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusForbidden but the handler returned this status code")
	// }

	// // SUCCESS CASE
	// rr = buildCtxSpecificUser(t, "PATCH", "application/json", correctNewUser, "me", "me", true, false)
	// if status := rr.Code; status == http.StatusForbidden {
	// 	t.Errorf(
	// 		"we did not expect a http.StatusForbidden but the handler returned this status code")
	// }

	// // User is authorized, but not allowed to access user id 123
	// // FAIL CASE
	// rr = buildCtxSpecificUser(t, "PATCH", "alication/json", correctNewUser, "123", "1234", true, true)
	// if status := rr.Code; status != http.StatusUnsupportedMediaType {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	// }

	// // malformed path
	// // FAIL CASE
	// rr = buildCtxSpecificUser(t, "PATCH", "alication/json", correctNewUser, "m", "1234", true, true)
	// if status := rr.Code; status != http.StatusUnsupportedMediaType {
	// 	t.Errorf(
	// 		"we expected an http.StatusUnsupportedMediaType but the handler returned wrong status code")
	// }

	// // Checking for correct headers
}

// TestSessionsHandler does something
// All tests pass! (7 total)
func TestSessionsHandler(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		method      string
		encoding    string
		nu          map[string]string
		pathExtras  string
		status      int
		expected    string
		expectedErr bool
	}{
		{
			"POST method header set correctly",
			"Make sure you're using the correct http method",
			"POST",
			"application/json",
			correctCreds,
			"",
			http.StatusMethodNotAllowed,
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code",
			false,
		},
		{
			"Character encoding header set correctly",
			"Make sure you're using the correct character encoding",
			"POST",
			"application/json",
			correctCreds,
			"",
			http.StatusUnsupportedMediaType,
			"we did not expect a http.StatusUnsupportedMediaType but the handler returned this status code",
			false,
		},
		{
			"User is authorized to begin a session",
			"Make sure there is an associate user with the given credentials",
			"POST",
			"application/json",
			correctCreds,
			"",
			http.StatusUnauthorized,
			"we did not expect a http.StatusUnauthorized but the handler returned this status code",
			false,
		},
	}
	failcases := []struct {
		name        string
		hint        string
		method      string
		encoding    string
		nu          map[string]string
		pathExtras  string
		status      int
		expected    string
		expectedErr bool
	}{
		{
			"PATCH method header not supported by this method",
			"Make sure you're using the correct http method",
			"PATCH",
			"application/json",
			incorrectEmailCreds,
			"",
			http.StatusMethodNotAllowed,
			"we expected a http.StatusMethodNotAllowed but the handler did not return this status code",
			false,
		},
		{
			"Character encoding header set incorrectly",
			"Make sure you're using the correct character encoding",
			"POST",
			"alication/json",
			incorrectEmailCreds,
			"",
			http.StatusUnsupportedMediaType,
			"we expected a http.StatuStatusUnsupportedMediaTypesMethodNotAllowed but the handler did not return this status code",
			false,
		},
		{
			"Character encoding header set incorrectly",
			"Make sure you're using the correct character encoding",
			"POST",
			"application/json",
			incorrectEmailCreds,
			"",
			http.StatusUnauthorized,
			"we expected a http.StatusUnauthorized but the handler did not return this status code",
			true,
		},
		{
			"Character encoding header set incorrectly",
			"Make sure you're using the correct character encoding",
			"POST",
			"application/json",
			incorrectPassCreds,
			"",
			http.StatusUnauthorized,
			"we expected a http.StatusUnauthorized but the handler did not return this status code",
			true,
		},
	}
	for i, c := range cases {
		// SUCCESS CASE
		rr := buildCtxSession(t, c.method, c.encoding, c.nu, c.pathExtras, c.expectedErr)
		if status := rr.Code; status == c.status {
			t.Errorf(c.expected)
		}
		// FAIL CASE
		fc := failcases[i]
		rr = buildCtxSession(t, fc.method, fc.encoding, fc.nu, fc.pathExtras, fc.expectedErr)
		if status := rr.Code; status != fc.status {
			t.Errorf(fc.expected)
		}
	}
	fc := failcases[len(failcases)-1]
	rr := buildCtxSession(t, fc.method, fc.encoding, fc.nu, fc.pathExtras, fc.expectedErr)
	if status := rr.Code; status != fc.status {
		t.Errorf(fc.expected)
	}
}

// TestSpecificSessionsHandler does something
// EndSession test cases (2) not operational of 6 total
func TestSpecificSessionsHandler(t *testing.T) {
	cases := []struct {
		name       string
		hint       string
		method     string
		encoding   string
		nu         map[string]string
		pathExtras string
		status     int
		expected   string
	}{
		{
			"DELETE method header set correctly",
			"Make sure you're using the correct http method",
			"DELETE",
			"application/json",
			correctNewUser,
			"",
			http.StatusMethodNotAllowed,
			"we did not expect a http.StatusMethodNotAllowed but the handler returned this status code",
		},
		{
			"Header encoding set correctly",
			"Make sure you're using the correct character encoding",
			"DELETE",
			"application/json",
			correctNewUser,
			"mine",
			http.StatusForbidden,
			"we did not expect a http.StatusForbidden but the handler returned this status code",
		},
		{
			"Session ending correctly (EndSession Test)",
			"Make sure you're using an existing session id",
			"DELETE",
			"application/json",
			correctNewUser,
			"mine",
			http.StatusInternalServerError,
			"we did not expect a http.StatusInternalServerError but the handler returned this status code",
		},
	}
	failcases := []struct {
		name       string
		hint       string
		method     string
		encoding   string
		nu         map[string]string
		pathExtras string
		status     int
		expected   string
	}{
		{
			"DELETE method header set incorrectly",
			"Make sure you're using the correct http method",
			"PATCH",
			"application/json",
			correctNewUser,
			"",
			http.StatusMethodNotAllowed,
			"we expected a http.StatusMethodNotAllowed but the handler did not return this status code",
		},
		{
			"Header encoding set incorrectly",
			"Make sure you're using the correct encoding scheme",
			"DELETE",
			"alication/json",
			correctNewUser,
			"",
			http.StatusForbidden,
			"we expected a http.StatusForbidden but the handler did not return this status code",
		},
		{
			"Pass a signing key that does not exist in sessions (EndSession Fails)",
			"Make sure you're passing an existing session id",
			"DELETE",
			"alication/json",
			correctNewUser,
			"",
			http.StatusInternalServerError,
			"we expected a http.StatusInternalServerError but the handler did not return this status code",
		},
	}
	for i, c := range cases {
		// SUCCESS CASE
		rr := buildCtxSpecificSession(t, c.method, c.encoding, c.nu, c.pathExtras)
		if status := rr.Code; status == c.status {
			t.Errorf(c.expected)
		}
		// FAIL CASE
		fc := failcases[i]
		rr = buildCtxSpecificSession(t, fc.method, fc.encoding, fc.nu, fc.pathExtras)
		if status := rr.Code; status != fc.status {
			t.Errorf(fc.expected)
		}
	}
}
