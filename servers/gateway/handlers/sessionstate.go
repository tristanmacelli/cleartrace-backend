package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"time"
)

// SessionState is a struct
type SessionState struct {
	BeginTime time.Time
	User      *users.User
}
