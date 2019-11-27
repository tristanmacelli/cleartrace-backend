package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Notify A simple store to store all the connections
type Notify struct {
	// Connections map[string]*websocket.Conn
	Connections []*websocket.Conn
	lock        sync.Mutex
}

// NewNotify does something
func NewNotify(connections []*websocket.Conn, lock sync.Mutex) *Notify {
	return &Notify{connections, lock}
}

// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// Thread-safe method for inserting a connection
func (ctx *HandlerContext) InsertConnection(conn *websocket.Conn) int {
	s := ctx.SocketStore
	s.lock.Lock()
	connID := len(s.Connections)
	// insert socket connection
	s.Connections = append(s.Connections, conn)
	s.lock.Unlock()
	return connID
}

// Thread-safe method for inserting a connection
func (ctx *HandlerContext) RemoveConnection(connID int) {
	s := ctx.SocketStore
	s.lock.Lock()
	// insert socket connection
	s.Connections = append(s.Connections[:connID], s.Connections[connID+1:]...)
	s.lock.Unlock()
}

// Simple method for writing a message to all live connections.
// In your homework, you will be writing a message to a subset of connections
// (if the message is intended for a private channel), or to all of them (if the message
// is posted on a public channel
func (ctx *HandlerContext) WriteToAllConnections(messageType int, data []byte) error {
	s := ctx.SocketStore

	var writeError error

	for _, conn := range s.Connections {
		writeError = conn.WriteMessage(messageType, data)
		if writeError != nil {
			return writeError
		}
	}

	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//TODO : DO WE REALLY NEED THIS???
	// CheckOrigin: func(r *http.Request) bool {
	// 	// This function's purpose is to reject websocket upgrade requests if the
	// 	// origin of the websockete handshake request is coming from unknown domains.
	// 	// This prevents some random domain from opening up a socket with your server.
	// 	return r.Header.Get("Origin") == "https://a2.sauravkharb.me"
	// },
}

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

// WebSocketConnectionHandler does something
func (ctx *HandlerContext) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// problem getting Session State
	// TODO: how do we handle ctx && socketStore as receivers

	if ctx.SessionStore == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Status Code 403: Unauthorized"))
		return
	}
	// handle the websocket handshake

	if r.Header.Get("Origin") != "https://a2.sauravkharb.me" {
		http.Error(w, "Websocket Connection Refused", 403)
		return
	}

	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}

	connID := ctx.InsertConnection(conn)
	// Invoke a goroutine for handling control messages from this connection
	go (func(conn *websocket.Conn, connID int) {
		defer conn.Close()
		defer ctx.RemoveConnection(connID)
		ctx.echo(conn)
	})(conn, connID)

}

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

// echo does something
func (ctx *HandlerContext) echo(conn *websocket.Conn) {
	// for { // infinite loop
	// 	messageType, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println("Error reading message.", err)
	// 		conn.Close()
	// 		return
	// 	}
	// 	// fmt.Printf("Got message: %#v\n", p)
	// 	if err := conn.WriteMessage(messageType, p); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }
	// s := ctx.SocketStore

	for {
		messageType, p, err := conn.ReadMessage()

		if messageType == TextMessage || messageType == BinaryMessage {
			fmt.Printf("Client says %v\n", p)
			fmt.Printf("Writing %s to all sockets\n", string(p))

			// TODO : Make sure you are writing messages to only memebers of the private channel

			ctx.WriteToAllConnections(TextMessage, append([]byte("Got message: "), p...))
		} else if messageType == CloseMessage {
			fmt.Println("Close message received.")
			break
			// } else if messageType == PingMessage {
			// pingHandler := conn.PingHandler()
			// } else if messageType == PongMessage {
			// pongHandler := conn.PongHandler()
		} else if err != nil {
			fmt.Println("Error reading message.")
			break
		}
		// TA Question: Should we be ignoring ping and pong messages
		// Potential TODO: Handling a ping message sent by client when the client wants to know the server is still alive, and sending a pong message back

		// Potential TODO: Handling a pong message when the server sends the client a ping, and the client responds with a pong
	}
	// What should we be doing as a part of this cleanup
	// cleanup
}

// func main() {
// 	mux := http.NewServeMux()

// 	ctx := socketStore{
// 		Connections: []*websocket.Conn{},
// 	}

// 	mux.HandleFunc("/ws", ctx.webSocketConnectionHandler)
// 	log.Fatal(http.ListenAndServe(":4001", mux))
// }
