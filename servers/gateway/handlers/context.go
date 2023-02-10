package handlers

import (
	"server-side-mirror/servers/gateway/indexes"
	"server-side-mirror/servers/gateway/models/users"
	"server-side-mirror/servers/gateway/sessions"
)

// HandlerContext a handler context struct that
// is a receiver on any of your HTTP
// handler functions that need access to
// globals, such as the key used for signing
// and verifying SessionIDs, the session store
// and the user store
type HandlerContext struct {
	Key          string
	UserStore    users.Store
	UserIndexes  indexes.Trie
	SessionStore sessions.Store
	SocketStore  Notify
}

// NewHandlerContext does something
func NewHandlerContext(key string, userStore users.Store, userIndexes indexes.Trie, sessionStore sessions.Store, socketStore Notify) *HandlerContext {
	if len(key) == 0 {
		panic("No User key")
	} else if userStore == nil {
		panic("No user")
	} else if &userIndexes == new(indexes.Trie) {
		panic("No Indexes")
	} else if sessionStore == nil {
		panic("No Session")
	} else if &socketStore == new(Notify) {
		panic("No Socket")
	}
	return &HandlerContext{key, userStore, userIndexes, sessionStore, socketStore}
}
