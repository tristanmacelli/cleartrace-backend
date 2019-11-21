package main

import (
	"assignments-Tristan6/servers/gateway/handlers"
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

// IndexHandler does something
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from API Gateway"))
}

type Director func(r *http.Request)

func CustomDirector(target *url.URL) Director {
	return func(r *http.Request) {
		r.Header.Add("X-User", r.Host)
		r.Host = target.Host
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
	}
}

//main is the main entry point for the server
func main() {
	address := os.Getenv("ADDR")
	// Default address the server should listen on
	if len(address) == 0 {
		address = ":443"
	}
	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	sessionkey := os.Getenv("SESSIONKEY")
	redisaddr := os.Getenv("REDISADDR")
	dsn := os.Getenv("DSN")

	messagesaddr := os.Getenv("MESSAGESADDR")
	messagesaddr1 := strings.Split(messagesaddr, ",")[0]
	// messagesaddr2 := strings.Split(messagesaddr, ",")[1]

	summaryaddr := os.Getenv("SUMMARYADDR")

	// create redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisaddr, // use default Addr
	})
	redisStore := sessions.NewRedisStore(redisClient, 0)

	userStore := users.NewMysqlStore(dsn)

	// If there are multiple addresses for either messages or summary then do the following
	// TODO: random number generator to pick between the available addresses
	u, err := url.Parse(messagesaddr1)
	if err != nil {
		fmt.Print(err)
	}
	messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(u)}

	// proxy
	// messagesProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: messagesaddr1})
	summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: summaryaddr})

	ctx := handlers.NewHandlerContext(sessionkey, userStore, redisStore)
	// starting a new mux session
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificUserHandler)
	mux.Handle("/v1/summary", summaryProxy)
	mux.Handle("/v1/channels", messagesProxy)
	mux.Handle("/v1/channels/{channelID}", messagesProxy)
	mux.Handle("/v1/channels/{channelID}/members", messagesProxy)
	mux.Handle("/v1/messages/{messageID}", messagesProxy)
	wrappedMux := handlers.NewLogger(mux)

	// logging server location or errors
	log.Printf("server is listening at %s...", address)
	log.Fatal(http.ListenAndServeTLS(address, tlsCertPath, tlsKeyPath, wrappedMux))

	/* To host server:
	- change path until in folder with main.go in it
	- 'go install main.go'
	- navigate to 'go' bin folder and run main.exe
	*/
}
