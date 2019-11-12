package handlers

import "net/http"

type Channel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
	Members     int    `json:"members,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	Creator     string `json:"creator,omitempty"`
	EditedAt    string `json:"editedAt,omitempty"`
}

//PageSummary represents summary properties for a web page
type Message struct {
	ID        string `json:"id,omitempty"`
	ChannelID string `json:"channelID,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Body      string `json:"body,omitempty"`
	Creator   string `json:"creator,omitempty"`
	EditedAt  string `json:"editedAt,omitempty"`
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if !correctHeader(w, r) {
		return
	}
}

func correctHeader(w http.ResponseWriter, r *http.Request) bool {
	// check for correct header
	ctype := r.Header.Get("X-User")
	if len(ctype) < 1 {
		// throw error
		http.Error(w, "must be authenticated", http.StatusUnauthorized)
		return false
	}
	return true
}
