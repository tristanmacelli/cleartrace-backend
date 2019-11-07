package users

import (
	"errors"
	"net/http"
)

var getByIDnextReturn = User{}
var getByEmailnextReturn = User{}
var getByUserNamenextReturn = User{}
var insertnextReturn = User{}
var updatenextReturn = User{}
var errNext = errors.New("generic error")

//MockStore represents a mock user store
type MockStore struct {
	// fields
}

// NewMockStore does something
func NewMockStore() *MockStore {

	return &MockStore{}
}

// SetGetByID does something
func SetGetByID(user *User) {
	getByIDnextReturn = *user
}

// SetGetByEmailnextReturn does something
func SetGetByEmailnextReturn(user *User) {
	getByEmailnextReturn = *user
}

// SetGetByUserNamenextReturn does something
func SetGetByUserNamenextReturn(user *User) {
	getByUserNamenextReturn = *user
}

// SetInsertnextReturn does something
func SetInsertnextReturn(user *User) {
	insertnextReturn = *user
}

// SetUpdatenextReturn does something
func SetUpdatenextReturn(user *User) {
	updatenextReturn = *user
}

// SetErr does something
func SetErr(err error) {
	errNext = err
}

//GetByID returns the User with the given ID
func GetByID(id int64) (*User, error) {
	return &getByIDnextReturn, errNext
}

//GetByEmail returns the User with the given email
func GetByEmail(email string) (*User, error) {
	return &getByEmailnextReturn, errNext
}

//GetByUserName returns the User with the given Username
func GetByUserName(username string) (*User, error) {
	return &getByUserNamenextReturn, errNext
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func Insert(user *User) (*User, error) {
	return &insertnextReturn, errNext
}

// LogSuccessfulSignIns does something
func LogSuccessfulSignIns(user *User, r *http.Request) {

}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func Update(id int64, updates *Updates) (*User, error) {
	return &updatenextReturn, errNext
}

//Delete deletes the user with the given ID
func Delete(id int64) error {
	return errNext
}
