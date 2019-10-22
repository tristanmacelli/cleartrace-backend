package users

import (
	"testing"
)

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
