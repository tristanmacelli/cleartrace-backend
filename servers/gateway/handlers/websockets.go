package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server-side-mirror/servers/gateway/sessions"
	"sync"

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
func (ctx *HandlerContext) WriteToAllConnections(curUserID int64, messageType int, data []byte) error {
	s := ctx.SocketStore
	var writeError error

	for _, conn := range s.Connections {
		// Either
		// reflect.DeepEqual(s.Connections[curUserID], conn)
		// OR cmp.Equal(s.Connections[curUserID], conn)
		writeError = conn.WriteMessage(messageType, data)
		if writeError != nil {
			return writeError
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
		w.Write([]byte("Status Code 401: Unauthorized"))
		return
	}
	// handle the websocket handshake

	if r.Header.Get("Origin") != "https://slack.tristanmacelli.com" {
		http.Error(w, "Websocket Connection Refused", 403)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}
	sessionState := &SessionState{}
	// Should we be using the sessionID or the userID when mapping connections?
	sessions.GetState(r, ctx.Key, ctx.SessionStore, sessionState)

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

// echo does something
func (ctx *HandlerContext) echo(conn *websocket.Conn) {
	connMQ, err := amqp.Dial("amqp://userMessageQueue")
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
		fmt.Println("Now actively listening for messages!")
		for d := range msgs {
			fmt.Println("Looping through infinitely!")
			fmt.Println("Message object is:", d)
			// THIS STATEMENT WILL BLOCK UNTIL ANOTHER MSG ARRIVES AT THE SRVR
			// AKA ALL STATEMENTS AFTER ReadJSON WILL NOT EXECUTE UNTIL NEW MSG
			// err := conn.ReadJSON(&d.Body)
			// if err != nil {
			// 	fmt.Println("Error reading JSON: ", err)
			// 	break
			// }
			// CURRENTLY THE CONTENT TYPE IS NIL
			// fmt.Println()
			// fmt.Println()
			// fmt.Println("Content Type of MQ message is: ", d.ContentType)
			// if d.ContentType == "application/json" {
			// d.Body is a uint8 type which is equivalent to byte
			fmt.Println("Handling Client Bound Messages")
			err := ctx.handleClientBoundMessages(d)
			if err != nil {
				// This should cause the connection to close
				fmt.Println("Error handling Client-bound messages: ", err)
				break
			}
			// } else {
			// 	// Close connection ?
			// }
			// TODO: Handle connection closure
			// } else if messageType == CloseMessage {
			// 	fmt.Println("Close message received.")
			// 	break
			// } else if err != nil {
			// 	fmt.Println("Error reading message.")
			// 	break
			// }
			// if err := d.Ack(false); err != nil {
			// 	log.Printf("Error acknowledging message : %s", err)
			// } else {
			// 	log.Printf("Acknowledged message")
			// 	// The Ack call above responds to rabbitMQ saying that we have received
			// 	// the message it sent us.
			// }
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	// What should we be doing as a part of this cleanup
	// cleanup
	fmt.Println("infinite loop exited")
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
	curUserID := message.Message.Creator.ID

	// Close connection (returning an error causes the for-loop above to break)
	if message.MessageType == "close-connection" {
		return errors.New("Closing connection")
	}

	// Broadcast to all clients when userIDs are not provided
	if userIDs == nil {
		log.Printf("Writing to all connections")
		err = ctx.WriteToAllConnections(curUserID, 1, d.Body)
		log.Printf("Wrote to all connections")
		if err != nil {
			log.Printf("Error writing message to all connections: %s", err)
			return err
		}
		// Broadcast to specific list when userIDs are provided
	} else {
		log.Printf("Writing to specific connections")
		err = ctx.WriteToSpecificConnections(curUserID, 1, d.Body, userIDs)
		log.Printf("Wrote to specific connections")
		if err != nil {
			log.Printf("Error writing message to specific connections: %s", err)
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
