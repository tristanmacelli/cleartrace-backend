package handlers

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

// getMessages gets the messages from the message queue
func getMessages() {
	conn, err := amqp.Dial("amqp://guest:guest@messagequeue:5672/")
	failOnError("Failed to open connection to RabbitMQ", err)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError("Failed to Open Channel", err)
	defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"helloQueue", // name
	// 	false,        // durable (do my messages last until I delete my connection)
	// 	false,        // delete when unused
	// 	false,        // exclusive
	// 	false,        // no-wait
	// 	nil,          // additional arguments
	// )
	// failOnError("Failed to declare queue", err)

	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	true,   // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// failOnError("Failed to register a consumer", err)

	// forever := make(chan bool)

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)
	// 	}
	// }()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever

}
