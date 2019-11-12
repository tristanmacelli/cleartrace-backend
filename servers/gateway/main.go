package main

import (
	"log"
	"net/http"
	"os"

	"./handlers"
)

//main is the main entry point for the server
func main() {
	address := os.Getenv("ADDR")
	// Default address the server should listen on
	if len(address) == 0 {
		address = ":80"
	}
	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	sessionkey := os.Getenv("SESSIONKEY")
	redisaddr := os.Getenv("REDISADDR")
	dsn := os.Getenv("DSN")
	// starting a new mux session
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", handlers.UserHandler)
	mux.HandleFunc("/v1/users/", handlers.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", handlers.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", handlers.SpecificUserHandler)
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	wrappedMux := NewLogger(mux)

	// logging server location or errors
	log.Printf("server is listening at %s...", address)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))

	/* To host server:
	- change path until in folder with main.go in it
	- 'go install main.go'
	- navigate to 'go' bin folder and run main.exe
	*/
}
