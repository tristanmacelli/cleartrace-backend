package users

import (
	"errors"
	"net/http"
	"server-side-mirror/servers/gateway/indexes"
)

var getByIDnextReturn = User{}
var getByEmailnextReturn = User{}
var getByUserNamenextReturn = User{}
var insertnextReturn = User{}
var updatenextReturn = User{}
var errNext = errors.New("generic error")

//MockStore represents a mock user store
type MockStore struct {
	// pointless string
}

// NewMockStore does something
func NewMockStore() *MockStore {
	return &MockStore{}
}

// SetGetByIDNextReturn does something
func SetGetByIDNextReturn(user *User) {
	getByIDnextReturn = *user
}

// SetGetByEmailNextReturn does something
func SetGetByEmailNextReturn(user *User) {
	getByEmailnextReturn = *user
}

// SetGetByUserNameNextReturn does something
func SetGetByUserNameNextReturn(user *User) {
	getByUserNamenextReturn = *user
}

// SetInsertNextReturn does something
func SetInsertNextReturn(user *User) {
	insertnextReturn = *user
}

// SetUpdateNextReturn does something
func SetUpdateNextReturn(user *User) {
	updatenextReturn = *user
}

// SetErr does something
func SetErr(err error) {
	errNext = err
}

//GetByID returns the User with the given ID
func (ms *MockStore) GetByID(id int64) (*User, error) {
	return &getByIDnextReturn, errNext
}

//GetByEmail returns the User with the given email
func (ms *MockStore) GetByEmail(email string) (*User, error) {
	return &getByEmailnextReturn, errNext
}

//GetByUserName returns the User with the given Username
func (ms *MockStore) GetByUserName(username string) (*User, error) {
	return &getByUserNamenextReturn, errNext
}

func (ms *MockStore) GetByIDs(ids []int64, orderBy string) (*[]User, error) {
	return &[]User{}, nil
}

func (ms *MockStore) IndexUsers(trie *indexes.Trie) {}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ms *MockStore) Insert(user *User) (*User, error) {
	return &insertnextReturn, errNext
}

// LogSuccessfulSignIns does something
func (ms *MockStore) LogSuccessfulSignIns(user *User, r *http.Request) {

}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *MockStore) Update(id int64, updates *Updates) (*User, error) {
	return &updatenextReturn, errNext
}

//Delete deletes the user with the given ID
func (ms *MockStore) Delete(id int64) error {
	return errNext
}
