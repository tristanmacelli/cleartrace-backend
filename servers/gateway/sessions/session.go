package sessions

import (
	"errors"
	"net/http"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	// Creating a new SessionID
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	// Saving the sessionState to the store
	err = store.Save(sessionID, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}
	// Adding a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)
	authValue := schemeBearer + sessionID.String()
	w.Header().Add(headerAuthorization, authValue)
	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	// Get the value of the Authorization header,
	id := r.Header.Get(headerAuthorization)
	// Or the "auth" query string parameter if no Authorization header is present,
	if id == "" {
		id = r.Header.Get(paramAuthorization)
	}
	// Then we validate the sessionID.
	sessionID, err := ValidateID(id, signingKey)
	// If the sessionID is not valid return the validation error.
	if err != nil {
		return InvalidSessionID, err
	}
	// If it's valid, return the SessionID.
	return sessionID, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	// Getting the SessionID from the request
	id, err := GetSessionID(r, signingKey)
	// If the sessionID is not valid return the validation error.
	if err != nil {
		return InvalidSessionID, err
	}
	// Getting the data associated with the returned SessionID from the given store
	err = store.Get(id, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}
	return id, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	// Getting the SessionID from the request
	id, err := GetSessionID(r, signingKey)
	// If the sessionID is not valid return the validation error.
	if err != nil {
		return InvalidSessionID, err
	}
	// Deleting the data associated with the returned sessionID from the given store
	store.Delete(id)
	return id, nil
}
