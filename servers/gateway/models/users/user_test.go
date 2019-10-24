package users

import (
	"testing"
)

// @saurav
// this script tests all the functions written in user.go
// same testing strategies used as other tests

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

// Test the (nu *NewUser) Validate() function to ensure it catches all possible validation errors,
// and returns no error when the new user is valid.

func TestValidateNewUser(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		nu          *NewUser
		expectError bool
	}{
		{
			"Invalid Email",
			"This is an invalid email so this should fail.",
			&NewUser{
				Email:        "#@%^%#$@#$@#.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword123",
				UserName:     "TMcGee123",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			true,
		},
		{
			"Password too short",
			"The password is too short but there was no recieved error.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "123",
				PasswordConf: "mypassword123",
				UserName:     "TMcGee123",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			true,
		},
		{
			"Passwords do not match",
			"The passwords didn't match but there was no recieved error.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword987",
				UserName:     "TMcGee123",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			true,
		},
		{
			"Invalid username: length 0",
			"The username is of length 0 but there was no recieved error.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword123",
				UserName:     "",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			true,
		},
		{
			"Invalid username: white space",
			"The username contained whie space but there was no recieved error.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword123",
				UserName:     "TMcGee 123",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			true,
		},
		{
			"Valid new user",
			"Things should pass here but there was a reiceved error.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword123",
				UserName:     "TMcGee123",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			false,
		},
	}

	for _, c := range cases {
		err := c.nu.Validate()
		if err != nil && !c.expectError { // There WAS an error but we DIDN'T expect one.
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectError && err == nil { // We DID expect and error but we DIDN'T recieve one.
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}

}

func TestToUser(t *testing.T) {
	cases := []struct {
		name        string
		hint        string
		nu          *NewUser
		expectError bool
	}{
		{
			"New user isn't validated",
			"Things shouldn't pass here but there was no error reiceved.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword123",
				UserName:     "",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			true,
		},
		{
			"Everything works",
			"Things should pass here but there was a reiceved error.",
			&NewUser{
				Email:        "myexampleEmail@live.com",
				Password:     "mypassword123",
				PasswordConf: "mypassword123",
				UserName:     "TMcGee123",
				FirstName:    "Tester",
				LastName:     "McGee",
			},
			false,
		},
	}
	for _, c := range cases {
		u, err := c.nu.ToUser()
		if u == nil && err != nil && !c.expectError { // There WAS an error but we DIDN'T expect one.
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
		}
		if u != nil && err == nil && c.expectError { // We DID expect and error but we DIDN'T recieve one.
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name     string
		hint     string
		u        *User
		expected string
	}{
		{
			"Missing First and Last names",
			"Make sure you're returning an empty string for missing name",
			&User{
				FirstName: "",
				LastName:  "",
			},
			"",
		},
		{
			"Missing First name only",
			"Make sure you're returning the last name even if the first name is empty",
			&User{
				FirstName: "",
				LastName:  "Something",
			},
			"Something",
		},
		{
			"Missing last name only",
			"Make sure you're returning the first name even if the last name is empty",
			&User{
				FirstName: "lastie",
				LastName:  "",
			},
			"lastie",
		},
		{
			"Both names present",
			"Make sure you're returning first name <space> last name",
			&User{
				FirstName: "firstie",
				LastName:  "lastie",
			},
			"firstie lastie",
		},
	}
	for _, c := range cases {
		str := c.u.FullName()
		if str != c.expected {
			t.Errorf("case %s: unexpected result \nHINT: %s", c.name, c.hint)
		}
	}

}

// cases
// could not crypt password (too long)
// normal password --> check if length of u.passHash us not zero after test
func TestSetPassword(t *testing.T) {
	cases := []struct {
		name    string
		hint    string
		u       *User
		isError bool
	}{
		{
			"Healthy password",
			"Make sure you're successfully encrypting and setting the password hash to the user",
			&User{
				PassHash: []byte{},
			},
			false,
		},
	}

	for _, c := range cases {
		err := c.u.SetPassword("mypassword")
		if !c.isError && err == nil {
			// t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
		if len(c.u.PassHash) == 0 {
			t.Errorf("case %s: unexpected result \nHINT: %s", c.name, c.hint)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		name    string
		hint    string
		u       *User
		ptPass  string
		isError bool
	}{
		{
			"Password and hash should match",
			"Make sure you authenticate correctly if the param poassword and stored hash match",
			&User{
				PassHash: []byte{},
			},
			"hello",
			false,
		},
		{
			"Password and hash shouldn't match",
			"Make sure you retiurn an error if the param poassword and stored hash don't match",
			&User{
				PassHash: []byte{},
			},
			"hello",
			true,
		},
	}
	for _, c := range cases {
		if c.name == "Password and hash should match" {
			_ = c.u.SetPassword(c.ptPass)
			err := c.u.Authenticate(c.ptPass)
			if !c.isError && err != nil {
				// t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
			}
		} else {
			_ = c.u.SetPassword("not hello")
			err := c.u.Authenticate(c.ptPass)
			if c.isError && err == nil {
				t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
			}
		}

	}

}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name    string
		hint    string
		up      *Updates
		isError bool
	}{
		{
			"The firstName must contain characters",
			"Make sure the firstName contains characters",
			&Updates{
				FirstName: "",
				LastName:  "lastname",
			},
			true,
		},
		{
			"The firstName must not contain white space",
			"Make sure the firstName contains white space",
			&Updates{
				FirstName: "firstname ",
				LastName:  "lastname",
			},
			true,
		},
		{
			"The lastName must contain characters",
			"Make sure the lastName contains characters",
			&Updates{
				FirstName: "firstname",
				LastName:  "",
			},
			true,
		},
		{
			"The lastName must not contain white space",
			"Make sure the lastName contains white space",
			&Updates{
				FirstName: "firstname",
				LastName:  "lastname ",
			},
			true,
		},
		{
			"The firstName and lastName are valid",
			"Everything should be working, but you received an error",
			&Updates{
				FirstName: "firstname",
				LastName:  "lastname",
			},
			false,
		},
	}
	var u User
	u.FirstName = "blank"
	u.LastName = "blank"
	for _, c := range cases {
		err := u.ApplyUpdates(c.up)
		if c.isError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
		if !c.isError && err != nil {
			t.Errorf("case %s: we didn't expected error but did get one\nHINT: %s", c.name, c.hint)
		}
	}
}
