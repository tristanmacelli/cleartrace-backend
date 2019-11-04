package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!
type SessionState struct {
	BeginTime time.Time
	User      users.User
}
