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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError("Failed to open connection to RabbitMQ", err)
	defer conn.Close()

	ch, err := conn.Channel()
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

	body := "Hello World"
	err = ch.Publish(
		"", //exchange
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError("Failed to publish message", err)
}
