package handlers

import "net/http"

/*
	A CORS middleware handler!
	See https://drstearns.github.io/tutorials/cors/ for help
*/

// Passer does something
type Passer struct {
	handler http.Handler
}

// NewLogger does something
func NewLogger(handlerToWrap http.Handler) *Passer {
	return &Passer{handlerToWrap}
}

// ServeHTTP does something
func (p *Passer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")
	p.handler.ServeHTTP(w, r)
}
