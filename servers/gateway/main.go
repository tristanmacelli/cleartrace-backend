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
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
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

func getAllUrls(addresses string) []*url.URL {
	urlSlice := strings.Split(addresses, ",")
	var urls []*url.URL
	for _, u := range urlSlice {
		url := url.URL{Scheme: "http", Host: u}
		urls = append(urls, &url)
	}
	return urls
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

	// create redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisaddr, // use default Addr
	})
	redisStore := sessions.NewRedisStore(redisClient, 0)
	dsn = fmt.Sprintf("root:%s@tcp("+dsn+")/users", os.Getenv("MYSQL_ROOT_PASSWORD"))
	userStore := users.NewMysqlStore(dsn)

	messagesaddr := os.Getenv("MESSAGEADDR")
	summaryaddr := os.Getenv("SUMMARYADDR")
	messagingUrls := getAllUrls(messagesaddr)
	// summaryUrls := getAllUrls(summaryaddr)

	// proxies
	messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(messagingUrls)}
	// summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(summaryUrls)}
	summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: summaryaddr})

	var conns []*websocket.Conn
	socketStore := handlers.NewNotify(conns, &sync.Mutex{})
	ctx := handlers.NewHandlerContext(sessionkey, userStore, redisStore, *socketStore)
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
	mux.HandleFunc("/v1/ws", ctx.WebSocketConnectionHandler)

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
