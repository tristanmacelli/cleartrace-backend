package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"server-side-mirror/servers/gateway/handlers"
	"server-side-mirror/servers/gateway/indexes"
	"server-side-mirror/servers/gateway/models/users"
	"server-side-mirror/servers/gateway/sessions"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// IndexHandler does something
func IndexHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello from API Gateway"))
}

// Director is a function wrapper
type Director func(request *http.Request)

// CustomDirector does load balancing using the round-robin method
func CustomDirector(targets []*url.URL, ctx *handlers.HandlerContext) Director {
	var counter int32 = 0

	return func(request *http.Request) {
		state := &handlers.SessionState{}
		_, err := sessions.GetState(request, ctx.Key, ctx.SessionStore, state)
		if err != nil {
			request.Header.Del("X-User")
			log.Println("Error getting User from GetState")
			log.Println(err)
			return
		}

		userJSON, _ := json.Marshal(state.User)
		userString := string(userJSON)

		targ := targets[counter%int32(len(targets))]
		atomic.AddInt32(&counter, 1)
		request.Header.Add("X-User", userString)
		request.Host = targ.Host
		request.URL.Host = targ.Host
		request.URL.Scheme = targ.Scheme
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

// main is the main entry point for the server
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
	messagingUrls := getAllUrls(messagesaddr)

	conns := make(map[int64]*websocket.Conn)
	socketStore := handlers.NewNotify(conns, &sync.Mutex{})
	indexedUsers := indexes.NewTrie(&sync.Mutex{})
	userStore.IndexUsers(indexedUsers)

	ctx := handlers.NewHandlerContext(sessionkey, userStore, *indexedUsers, redisStore, *socketStore)

	// proxies
	messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(messagingUrls, ctx)}

	mux := mux.NewRouter()
	mux.HandleFunc("/", IndexHandler)

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/{userID}", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/users/email/{email}", ctx.GetUserByEmailHandler)
	mux.HandleFunc("/v1/users/search/", ctx.SearchHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/mine", ctx.SpecificSessionsHandler)
	mux.HandleFunc("/v1/ws", ctx.WebSocketConnectionHandler)
	mux.Handle("/v1/channels/{channelID}/members", messagesProxy)
	mux.Handle("/v1/channels/{channelID}", messagesProxy)
	mux.Handle("/v1/channels", messagesProxy)
	mux.Handle("/v1/messages/{messageID}", messagesProxy)

	wrappedMux := handlers.NewLogger(mux)

	// logging server location or errors
	log.Printf("server is listening %s...", address)

	// TLS 1.3 configuration
	config := &tls.Config{
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			// 1.3 Cipher suites
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		PreferServerCipherSuites: true,
		ClientSessionCache:       tls.NewLRUClientSessionCache(128),
	}
	server := &http.Server{Addr: address, Handler: wrappedMux, TLSConfig: config}
	// ListenAndServeTLS serves http2 && tls 1.2 by default, we are using it to only use TLS1.3
	log.Fatal(server.ListenAndServeTLS(tlsCertPath, tlsKeyPath))

	// TODO: Fix CORS issue
	// log.Fatal(http3.ListenAndServeQUIC(address, tlsCertPath, tlsKeyPath, wrappedMux))

	// TLS1.2 (Works without much issue)
	// TODO: Remove eventually
	// log.Fatal(http.ListenAndServeTLS(address, tlsCertPath, tlsKeyPath, wrappedMux))

	/* To host server:
	- change path until in folder with main.go in it
	- 'go install main.go'
	- navigate to 'go' bin folder and run main.exe
	*/
}
