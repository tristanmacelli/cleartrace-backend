package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
)

//HandlerContext a handler context struct that
//is a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
type HandlerContext struct {
	Key          string
	UserStore    users.Store
	SessionStore sessions.Store
	SocketStore  Notify
}

// NewHandlerContext does something
func NewHandlerContext(key string, userStore users.Store, sessionStore sessions.Store, socketStore Notify) *HandlerContext {
	if len(key) == 0 {
		panic("No User key")
	} else if userStore == nil {
		panic("No user")
	} else if sessionStore == nil {
		panic("No Session")
	}
	return &HandlerContext{key, userStore, sessionStore, socketStore}
}
