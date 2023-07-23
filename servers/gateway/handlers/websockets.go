package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server-side-mirror/servers/gateway/sessions"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	// "github.com/streadway/amqp" // this package is deprecated
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
	ConnectionPool map[int64]*websocket.Conn
	lock           *sync.Mutex
}

type MessagingTransaction struct {
	Type      string  `json:"type"`
	UserIDs   []int64 `json:"userIDs"`
	Channel   Channel `json:"channel"`
	Message   Message `json:"message"`
	ChannelID string  `json:"channelID"`
	MessageID string  `json:"messageID"`
}

// Channel is from our messaging service
type Channel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	Members     []int64 `json:"members"`
	CreatedAt   string  `json:"createdat"`
	Creator     Creator `json:"creator"`
	EditedAt    string  `json:"editedat"`
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

// NewNotify creates a WebSocket connection pool and mutex
func NewNotify(connections map[int64]*websocket.Conn, lock *sync.Mutex) *Notify {
	return &Notify{connections, lock}
}

// InsertConnection is a thread-safe method for inserting a connection
func (ctx *HandlerContext) InsertConnection(conn *websocket.Conn, userID int64) int {
	s := ctx.SocketStore
	s.lock.Lock()
	connID := len(s.ConnectionPool)
	// insert socket connection
	s.ConnectionPool[userID] = conn
	s.lock.Unlock()
	return connID
}

// RemoveConnection is a thread-safe method for inserting a connection
func (ctx *HandlerContext) RemoveConnection(connID int, userID int64) {
	s := ctx.SocketStore
	s.lock.Lock()
	// remove socket connection
	delete(s.ConnectionPool, userID)
	s.lock.Unlock()
}

// WriteToAllConnections is a simple method for writing a message to all live connections.
// In your homework, you will be writing a message to a subset of connections
// (if the message is intended for a private channel), or to all of them (if the message
// is posted on a public channel
func (ctx *HandlerContext) WriteToAllConnections(creatorUserID int64, messageType int, data []byte) error {
	s := ctx.SocketStore

	for id, conn := range s.ConnectionPool {
		// Skip the loop iteration associated with the message creator
		if id == creatorUserID {
			continue
		}
		writeError := conn.WriteMessage(messageType, data)
		if writeError != nil {
			return writeError
		}
	}

	return nil
}

// WriteToSpecificConnections writes to specific connections denoted by the userIDs attached to the
// message being returned from the message queue
func (ctx *HandlerContext) WriteToSpecificConnections(creatorUserID int64, messageType int, data []byte, ids []int64) error {
	s := ctx.SocketStore

	for _, id := range ids {
		// Skip the loop iteration associated with the message creator
		if id == creatorUserID {
			continue
		}
		// Skip loop iteration(s) associated with dead connections
		conn := s.ConnectionPool[id]
		if conn == nil {
			continue
		}
		writeError := conn.WriteMessage(messageType, data)
		if writeError != nil {
			return writeError
		}
	}

	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(request *http.Request) bool {
		return request.Header.Get("Origin") == "https://slack.tristanmacelli.com"
		// Uncomment this to accept WebSocket connections from localhost
		// return true
	},
}

// WebSocketConnectionHandler upgrades clients to a WebSocket connection
// and adds that connection to a connection pool. Connections in the connection
// pool can be notified when events are read from the RabbitMQ server
func (ctx *HandlerContext) WebSocketConnectionHandler(response http.ResponseWriter, request *http.Request) {
	// problem getting Session State
	// TODO: how do we handle ctx && socketStore as receivers

	if ctx.SessionStore == nil {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Status Code 401: Unauthorized"))
		return
	}

	// Comment this out to accept WebSocket connections from localhost
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

	// Channels (or chan) are an inter-process communication (IPC) construct allowing for typed message passing between
	// different goroutines
	waitOnErrorOrClose := make(chan bool)
	exit := make(chan bool)

	go func() {
		log.Println("Listening for messages from the Queue!")
		for {
			// select statements allow a goroutine to read from multiple channels at once
			// Reading from multiple channels simultaneously in Golang:
			// (source: https://stackoverflow.com/questions/20593126/reading-from-multiple-channels-simultaneously-in-golang)
			// Reading from multiple channels and creating a second goroutine to read WS messages (blocking but not chan)
			// (source: https://stackoverflow.com/questions/6807590/how-to-stop-a-goroutine)
			select {
			case d := <-msgs:
				message := &MessagingTransaction{}
				if err := json.Unmarshal(d.Body, message); err != nil {
					failOnError("Error decoding message JSON", err)
					CloseClientConnection(conn)
					// close(chan_reference) frees any receiving goroutines to unblock
					// In this case we are allowing the echo function to end (it is currently blocked on the <-waitForErrorOrClose
					// line)
					close(waitOnErrorOrClose)
					// Return ends the goroutine (but does not alter the wrapping goroutine, only close(chan_reference)
					//  can unblock receiving goroutines)
					return
				}

				err := ctx.HandleOutboundMessages(d, message)
				if err != nil {
					// This should cause the connection to close
					failOnError("Error handling Client-bound messages", err)
					CloseClientConnection(conn)
					close(waitOnErrorOrClose)
					return
				}
			case <-exit:
				log.Println("An inbound message caused a connection closure")
				close(waitOnErrorOrClose)
				return
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	go func() {
		log.Println("Listening for messages from the Client!")
		for {
			messageType, msg, err := conn.ReadMessage() // This is a blocking statement
			if err != nil {
				// Log any error that isn't a normal connection close
				if err.Error() != "websocket: close 1000 (normal): user terminated session" {
					failOnError("Error reading messages from Client", err)
				}
				// If the queue listener goroutine exits first, the parent goroutine will exit and close the connection
				// causing a read error. In this case, the client connection has already been closed
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					CloseClientConnection(conn) // Close WebSocket client connection on error
				}

				exit <- true
				close(exit)
				break
			}
			err, exitGoRoutine := HandleInboundMessages(messageType, msg, conn)
			if err != nil {
				failOnError("Error handling messages from Client", err)
				exit <- true
				close(exit)
				break
			}
			// In case of a CloseMessage received from a client that was successfully handled
			// and no errors were produced we want to exit the goroutine
			if exitGoRoutine {
				log.Println(
					"Successfully sent a close message to the client. Ending the goroutine & closing the connection server side",
				)
				exit <- true
				close(exit)
				break
			}
		}
	}()

	// Receing from the waitOnErrorOrClose channel blocks this goroutine until the server receives a close message or
	// produces an error
	<-waitOnErrorOrClose
	log.Println("websocket client closed")
}

// HandleOutboundMessages forwards messages from the RabbitMQ message queue (messaging microservice) to the
// client connections that are linked to the message
func (ctx *HandlerContext) HandleOutboundMessages(d amqp.Delivery, message *MessagingTransaction) error {
	userIDs := message.UserIDs
	creatorUserID := message.Message.Creator.ID

	// Broadcast to all clients when userIDs are not provided
	if userIDs == nil {
		err := ctx.WriteToAllConnections(creatorUserID, 1, d.Body)
		if err != nil {
			log.Printf("Error writing message to all connections: %s", err)
			return err
		}
		// Broadcast to specific list when userIDs are provided
	} else {
		err := ctx.WriteToSpecificConnections(creatorUserID, 1, d.Body, userIDs)
		if err != nil {
			log.Printf("Error writing message to specific connections: %s", err)
			return err
		}
	}
	return nil
}

// HandleInboundMessages
func HandleInboundMessages(messageType int, msg []byte, conn *websocket.Conn) (error, bool) {
	if messageType == CloseMessage {
		return CloseClientConnection(conn)
	}
	return nil, false
}

func CloseClientConnection(conn *websocket.Conn) (error, bool) {
	log.Println("Closing current WebSocket connection")
	cm := websocket.FormatCloseMessage(
		websocket.CloseNormalClosure,
		"WebSocket connection closed cleanly",
	)
	if err := conn.WriteControl(websocket.CloseMessage, cm, time.Now().Add(writeWait)); err != nil {
		return err, true
	}
	return nil, true
}

func failOnError(msg string, err error) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}
