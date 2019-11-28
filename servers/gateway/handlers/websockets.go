package handlers

import (
	"assignments-Tristan6/servers/gateway/models/users"
	"assignments-Tristan6/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// Notify A simple store to store all the connections
type Notify struct {
	Connections map[int64]*websocket.Conn
	lock        *sync.Mutex
}

type mqMessage struct {
	MessageType string  `json:"type"`
	Channel     Channel `json:"channel"`
	Message     Message `json:"message"`
	UserIDs     []int64 `json:"userIDs"`
	ChannelID   string  `json:"channelID"`
	MessageID   string  `json:"messageID"`
}

// Channel is from our messaging service
type Channel struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Private     bool     `json:"private"`
	Members     []string `json:"members"`
	CreatedAt   string   `json:"createdat"`
	Creator     string   `json:"creator"`
	EditedAt    string   `json:"editedat"`
}

// Message is from our messaging service
type Message struct {
	ID        string `json:"_id"`
	ChannelID string `json:"channelid"`
	CreatedAt string `json:"createdat"`
	Body      string `json:"body"`
	Creator   string `json:"creator"`
	EditedAt  string `json:"editedat"`
}

// NewNotify does something
func NewNotify(connections map[int64]*websocket.Conn, lock *sync.Mutex) *Notify {
	return &Notify{connections, lock}
}

// InsertConnection is a thread-safe method for inserting a connection
func (ctx *HandlerContext) InsertConnection(conn *websocket.Conn, userID int64) int {
	s := ctx.SocketStore
	s.lock.Lock()
	connID := len(s.Connections)
	// insert socket connection
	// TODO: map userID to the associated web socket connection
	s.Connections[userID] = conn
	s.lock.Unlock()
	return connID
}

// RemoveConnection is a thread-safe method for inserting a connection
func (ctx *HandlerContext) RemoveConnection(connID int, userID int64) {
	s := ctx.SocketStore
	s.lock.Lock()
	// insert socket connection
	delete(s.Connections, userID)
	s.lock.Unlock()
}

// WriteToAllConnections is a simple method for writing a message to all live connections.
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

// WriteToSpecificConnections writes to specific connections denoted by the userIDs attached to the
// message being returned from the message queue
func (ctx *HandlerContext) WriteToSpecificConnections(messageType int, data []byte, ids []int64) error {
	s := ctx.SocketStore
	var writeError error

	for _, id := range ids {
		conn := s.Connections[id]
		if conn != nil {
			writeError = conn.WriteMessage(messageType, data)
			if writeError != nil {
				return writeError
			}
		}
	}

	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	sessionID := r.Header.Get("auth")
	// TODO: check that this works?
	var user *users.User
	var sessionState SessionState
	sessionState.User = user
	sessionState.BeginTime = time.Now()
	sessions.GetState(r, sessionID, ctx.SessionStore, sessionState) // pass in sessionState struct

	// TODO: access to set with the specific connection and pass it to this function
	connID := ctx.InsertConnection(conn, sessionState.User.ID) // pass in the id that is contained within the sessionState struct
	// Invoke a goroutine for handling control messages from this connection
	go (func(conn *websocket.Conn, connID int) {
		defer conn.Close()
		defer ctx.RemoveConnection(connID, sessionState.User.ID)
		ctx.echo(conn)
	})(conn, connID)

}

// echo starts a goroutine that connects to the RabbitMQ server, reads events off the queue,
// and broadcasts them to all of the existing WebSocket connections that should hear about
// that event. If you get an error writing to the WebSocket, just close it and remove it from
// the list (client went away without closing from their end). Also make sure you start a read
// pump that reads incoming control messages, as described in the Gorilla WebSocket API
// documentation: http://godoc.org/github.com/gorilla/websocket

// echo does something
func (ctx *HandlerContext) echo(conn *websocket.Conn) {
	connMQ, err := amqp.Dial("amqp://guest:guest@messagequeue:5672/")
	failOnError("Failed to open connection to RabbitMQ", err)
	defer connMQ.Close()

	ch, err := connMQ.Channel()
	failOnError("Failed to Open Channel", err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"helloQueue", // name
		false,        // durable (do my messages last until I delete my connection)
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // additional arguments
	)
	failOnError("Failed to declare queue", err)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError("Failed to register a consumer", err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			if d.ContentType == "application/json" {
				err := ctx.handleClientBoundMessages(d)
				if err != nil {
					break
				}
			} else {
				// Close connection ?
			}
			// TODO: Handle closure
			// } else if messageType == CloseMessage {
			// 	fmt.Println("Close message received.")
			// 	break
			// } else if err != nil {
			// 	fmt.Println("Error reading message.")
			// 	break
			// }
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	// What should we be doing as a part of this cleanup
	// cleanup
}

// handleClientBoundMessages sends messages to
func (ctx *HandlerContext) handleClientBoundMessages(d amqp.Delivery) error {
	message := &mqMessage{}
	err := json.Unmarshal(d.Body, message)
	if err != nil {
		log.Printf("Error decoding message JSON: %s", err)
		return err
	}
	userIDs := message.UserIDs

	// Broadcast to all clients assuming we dont have userIDs
	if userIDs == nil {
		err = ctx.WriteToAllConnections(1, append([]byte("Got message: ")))
		if err != nil {
			log.Printf("Error decoding message JSON: %s", err)
			return err
		}
		// Broadcast to specific list assuming we have userIDs
	} else {
		err = ctx.WriteToSpecificConnections(1, append([]byte("Got message: ")), userIDs)
		if err != nil {
			log.Printf("Error decoding message JSON: %s", err)
			return err
		}
	}
	return nil
}

func failOnError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}
