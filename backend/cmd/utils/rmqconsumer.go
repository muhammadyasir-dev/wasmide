package utils

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)
conn, err := amqp.Dial("amqp://guest:guest://localhost:5672")
failOnError(err, "failed to listen rabbitmqe")

defer conn.Close()

ch, err := conn.Channel()
failOnError(err, "failed to listen rabbitmqe")
defer ch.Close()
q, err := ch.QueueDeclare(
"hello",
false,
false,
false,
false,
nil,
	)

failOnError(err, "failed to make que")

}


func failOnError(error error, msg string){

}

