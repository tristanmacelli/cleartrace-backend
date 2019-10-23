package users

import (
	"crypto/md5"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	// Validating the new user according to these rules
	// Checking for a valid Email address field (see mail.ParseAddress)
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		//log.Fatal(e, err, "Error: Email address is not in a valid format: yourname@domain.com")
		return err
	}

	// Checking that Password is at least 6 characters
	if len(nu.Password) <= 5 {
		return fmt.Errorf("Error: Password is not 6 characters or more")
	}

	// Checking that Password and PasswordConf match
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Error: Password doesnt not match the confirmed password")
	}
	// Checking that UserName is of non-zero length and does not contain spaces
	if len(nu.UserName) < 1 || strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("Error: Usernames must have a non-zero length and must contain no spaces")
	}
	//use fmt.Errorf() to generate appropriate error messages if
	//the new user doesn't pass one of the validation rules
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password

	//TODO: call Validate() to validate the NewUser and
	//return any validation errors that may occur.
	err := nu.Validate()
	fmt.Println(err)
	if err != nil { // there was an error.
		return nil, err
	}
	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	var us User
	us.Email = nu.Email

	err1 := us.SetPassword(nu.Password)
	if err1 != nil {
		return nil, err1
	}
	us.UserName = nu.UserName
	us.FirstName = nu.FirstName
	us.LastName = nu.LastName

	// The ID field will be left as a zero-value; the Store
	// implementation will set that field to the DBMS-assigned
	// primary key value.
	// The following sets the PhotoURL field to the Gravatar PhotoURL
	// for the user's email address.
	// see https://en.gravatar.com/site/implement/hash/
	// and https://en.gravatar.com/site/implement/images/ for more information.

	// Create new hash with md5 for photo url
	var hash = strings.ToLower(strings.TrimSpace(nu.Email))
	hasher := md5.New()
	hasher.Write([]byte(hash))
	us.PhotoURL = gravatarBasePhotoURL + hash

	// User struct
	/*ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`*/

	return &us, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
func (u *User) FullName() string {
	// If both first and last name are missing, this returns an empty string
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	// FirstName is an empty string, no space is put between the names.
	if u.FirstName == "" {
		return u.LastName
	}
	// LastName is an empty string, no space is put between the names.
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	// convert from string to byte hash
	bytePass := []byte(password)

	// convert byte pass to hash
	encryptedPassword, err := bcrypt.GenerateFromPassword(bytePass, bcryptCost)

	if err != nil {
		return err
	}
	u.PassHash = encryptedPassword

	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt

	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return err
	}

	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//TODO: set the fields of `u` to the values of the related
	//field in the `updates` struct

	if len(updates.FirstName) < 1 || strings.Contains(updates.FirstName, " ") {
		return fmt.Errorf("error: Firstname must contain characters and must not contain white space")
	}
	if len(updates.LastName) < 1 || strings.Contains(updates.LastName, " ") {
		return fmt.Errorf("error: Firstname must contain characters and must not contain white space")
	}
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName

	return nil
}
