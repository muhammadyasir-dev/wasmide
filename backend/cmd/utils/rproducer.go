package utils

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func Rabbitmqueproducer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "problem while building que ")

	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "channel failed to create")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"machina ist alive",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "declare queue")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)

	defer cancel()

	body := "booty"
	err = PublishwithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "error while publishing message")
	log.Printf("[x] sent data successfully")
}

func failOnError(error error, msg string) {

	if err != nil {
		log.Fatalf("failed section on fail on error expression evaluation %s\n", msg)

	}

}
