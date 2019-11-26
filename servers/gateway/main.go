package main

import (
	"assignments-Tristan6/servers/gateway/handlers"
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"

	"github.com/go-redis/redis"
)

// IndexHandler does something
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from API Gateway"))
}

// Director is a function wrapper
type Director func(r *http.Request)

// CustomDirector does load balancing using the round-robin method
func CustomDirector(targets []*url.URL) Director {
	var counter int32
	counter = 0

	return func(r *http.Request) {
		targ := targets[counter%int32(len(targets))]
		atomic.AddInt32(&counter, 1)
		r.Header.Add("X-User", r.Host)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
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

	messagesaddr := os.Getenv("MESSAGEADDR")
	messagesaddrSlice := strings.Split(messagesaddr, ",")
	// messagesaddr1 := messagesaddrSlice[0]
	// messagesaddr2 := messagesaddrSlice[1]

	// u1 := url.URL{Scheme: "http", Host: messagesaddr1}
	// u2 := url.URL{Scheme: "http", Host: messagesaddr2}

	// urlSlice := []*url.URL{&u1, &u2}
	var urlSlice []*url.URL
	// var messagingUrls []*url.URL
	for _, u := range messagesaddrSlice {
		url := url.URL{Scheme: "http", Host: u}
		urlSlice = append(urlSlice, &url)
	}

	summaryaddr := os.Getenv("SUMMARYADDR")
	// summaryaddrSlice := strings.Split(summaryaddr, ",")

	// var summaryUrls []*url.URL
	// for _, u := range summaryaddrSlice {
	// 	url := url.URL{Scheme: "http", Host: u}
	// 	summaryUrls = append(urlSlice, &url)
	// }

	// create redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisaddr, // use default Addr
	})
	redisStore := sessions.NewRedisStore(redisClient, 0)
	userStore := users.NewMysqlStore(dsn)

	// proxies
	messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(urlSlice)}
	// summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(summaryUrls)}
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
	log.Printf("server is listening testing %s...", address)
	log.Fatal(http.ListenAndServeTLS(address, tlsCertPath, tlsKeyPath, wrappedMux))

	/* To host server:
	- change path until in folder with main.go in it
	- 'go install main.go'
	- navigate to 'go' bin folder and run main.exe
	*/
}
