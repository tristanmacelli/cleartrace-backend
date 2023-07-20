package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server-side-mirror/servers/gateway/sessions"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

const writeWait = time.Second

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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Private     bool     `json:"private"`
	Members     []string `json:"members"`
	CreatedAt   string   `json:"createdat"`
	Creator     Creator  `json:"creator"`
	EditedAt    string   `json:"editedat"`
}

// Message is from our messaging service
type Message struct {
	ID        string  `json:"id"`
	ChannelID string  `json:"channelid"`
	CreatedAt string  `json:"createdat"`
	Body      string  `json:"body"`
	Creator   Creator `json:"creator"`
	EditedAt  string  `json:"editedat"`
}

type Creator struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
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
func (ctx *HandlerContext) WriteToAllConnections(curUserID int64, messageType int, data []byte) error {
	s := ctx.SocketStore
	var writeError error

	for id, conn := range s.Connections {
		if id != curUserID {
			writeError = conn.WriteMessage(messageType, data)
			if writeError != nil {
				return writeError
			}
		}
	}

	return nil
}

// WriteToSpecificConnections writes to specific connections denoted by the userIDs attached to the
// message being returned from the message queue
func (ctx *HandlerContext) WriteToSpecificConnections(curUserID int64, messageType int, data []byte, ids []int64) error {
	s := ctx.SocketStore
	var writeError error

	for _, id := range ids {
		if id != curUserID {
			conn := s.Connections[id]
			if conn != nil {
				writeError = conn.WriteMessage(messageType, data)
				if writeError != nil {
					return writeError
				}
			}
		}
	}

	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketConnectionHandler upgrades clients to a WebSocket connection
// and adds that connection to a list of connections to notify when events
// are read from the RabbitMQ server
func (ctx *HandlerContext) WebSocketConnectionHandler(response http.ResponseWriter, request *http.Request) {
	// problem getting Session State
	// TODO: how do we handle ctx && socketStore as receivers

	if ctx.SessionStore == nil {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Status Code 401: Unauthorized"))
		return
	}
	// handle the websocket handshake

	if request.Header.Get("Origin") != "https://slack.tristanmacelli.com" {
		http.Error(response, "Websocket Connection Refused", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		http.Error(response, "Failed to open websocket connection", http.StatusUnauthorized)
		return
	}
	sessionState := &SessionState{}
	// Should we be using the sessionID or the userID when mapping connections?
	sessions.GetState(request, ctx.Key, ctx.SessionStore, sessionState)

	// Access to set with the specific connection and pass it to this function
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
func (ctx *HandlerContext) echo(conn *websocket.Conn) {
	connMQ, err := amqp.Dial("amqp://userMessageQueue")
	failOnError("Failed to open connection to RabbitMQ", err)
	defer connMQ.Close()

	ch, err := connMQ.Channel()
	failOnError("Failed to Open Channel", err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messageLoopbackQueue", // name (was previously helloQueue)
		false,                  // durable (do my messages last until I delete my connection)
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // additional arguments
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
		fmt.Println("Now actively listening for messages!")
		for {
			// TODO: Figure out how to safely unblock from this statement when sending a message
			// Potential Solution: 2 Go routines this one continues to handle messages from the q
			// and the 2nd one will specifically handle close messages sent by logging out, refreshing,
			// or closing the browser tab/window
			// Search: gorilla websocket read message without blocking
			// Read: https://github.com/gorilla/websocket/issues/81
			// fmt.Println("1")
			// Close connection (returning an error causes the for-loop above to break)
			// When this logic is called it works properly, but blocks in any other case
			// if msg, err := HandleConnectionClosure(conn); err != nil {
			// 	break
			// } else if msg != nil {
			// fmt.Printf("Client says %v\n", msg)
			// }
			// fmt.Println("2")
			d := <-msgs
			fmt.Println("Looping through infinitely!")

			message := &mqMessage{}
			err := json.Unmarshal(d.Body, message)
			if err != nil {
				log.Printf("Error decoding message JSON: %s", err)
				break
			}

			err = ctx.handleClientBoundMessages(d, message)
			if err != nil {
				// This should cause the connection to close
				fmt.Println("Error handling Client-bound messages: ", err)
				break
			}
			// fmt.Println("3")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	// What should we be doing as a part of this cleanup
	// cleanup
	fmt.Println("infinite loop exited")
}

// handleClientBoundMessages forwards messages from the messaging microservice to the
// clients that are linked to the message
func (ctx *HandlerContext) handleClientBoundMessages(d amqp.Delivery, message *mqMessage) error {
	userIDs := message.UserIDs
	curUserID := message.Message.Creator.ID

	// Broadcast to all clients when userIDs are not provided
	if userIDs == nil {
		err := ctx.WriteToAllConnections(curUserID, 1, d.Body)
		if err != nil {
			log.Printf("Error writing message to all connections: %s", err)
			return err
		}
		// Broadcast to specific list when userIDs are provided
	} else {
		err := ctx.WriteToSpecificConnections(curUserID, 1, d.Body, userIDs)
		if err != nil {
			log.Printf("Error writing message to specific connections: %s", err)
			return err
		}
	}
	return nil
}

func HandleConnectionClosure(conn *websocket.Conn) ([]byte, bool) {
	messageType, msg, err := conn.ReadMessage()
	if err != nil || messageType == CloseMessage {
		fmt.Println("Closing current WebSocket connection")
		cm := websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"WebSocket connection closed cleanly",
		)
		if err := conn.WriteControl(websocket.CloseMessage, cm, time.Now().Add(writeWait)); err != nil {
			fmt.Println("Error:", err)
		}
		return []byte{}, true
	}
	return msg, false
}

func failOnError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}
