package main

import (
	"net/http"
	"os"
	"log"
)

// SummaryHandler, temporarily shows an html element at localhost/v1/summary
func SummaryHandler (w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello world</h1>"))
}

//main is the main entry point for the server
func main() {
	address := os.Getenv("ADDR")
	if len(address) == 0 {
		address = ":80"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", SummaryHandler)
	log.Printf("server is listening at %s...", address)
	log.Fatal(http.ListenAndServe(address, mux))

	/* To host server:
		- change path until in folder with main.go in it
		- 'go install main.go'
		- navigate to 'go' bin folder and run main.exe
	*/


	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
		the server should listen on. If empty, default to ":80"				DONE
	- Create a new mux for the web server.												  DONE
	- Tell the mux to call your handlers.SummaryHandler function
		when the "/v1/summary" URL path is requested.									DONE
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.								DONE
	*/
}
