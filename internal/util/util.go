package util

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func LogStruct(v interface {}) {
	structJSON, err := json.Marshal(v)
	if err != nil {
		log.Printf("Order: %+v", v)
	} else {
		log.Printf("Order: %s", structJSON)
	}
}

func GenerateUUID() string {
	uuid, err := uuid.NewRandom()
	FailOnError(err, "Failed to generate UUID")
	return uuid.String()
}

func DeclareRabbitMQChannel() (*amqp.Connection, *amqp.Channel, chan amqp.Confirmation) {
	conn, err := amqp.Dial(RabbitMQURL)
	FailOnError(err, "Failed to create rabbitmq connection")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	err = ch.Confirm(false)
	FailOnError(err, "Failed to set confirm mode")
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	return conn, ch, confirms
}