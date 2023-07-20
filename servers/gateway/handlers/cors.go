package handlers

import (
	"net/http"
)

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
func (p *Passer) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	response.Header().Set("Access-Control-Expose-Headers", "Authorization")
	response.Header().Set("Access-Control-Max-Age", "600")

	if request.Method == "OPTIONS" {
		response.WriteHeader(http.StatusOK)
		return
	}

	p.handler.ServeHTTP(response, request)

}
