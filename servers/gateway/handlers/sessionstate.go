package handlers

import (
	"server-side-mirror/servers/gateway/models/users"
	"time"
)

// SessionState is a struct
type SessionState struct {
	BeginTime time.Time   `json:"beginTime"`
	User      *users.User `json:"user"`
}
