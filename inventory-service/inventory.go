package inventoryservice

import (
	"fmt"
	"log"

	"github.com/gayaldassanayake/rabbitmq-order-system/internal/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RunService() {
	log.Printf("Inventory service is up and running")
	conn, ch, _ := util.DeclareRabbitMQChannel()
	defer conn.Close()
	defer ch.Close()
	// TODO: verify confirms

	err := ch.ExchangeDeclare(
		util.OrderEventsExchange,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err, fmt.Sprintf("Failed to declare exchange: %s", util.OrderEventsExchange))

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, fmt.Sprintf("Failed to declare queue: %s", q.Name))

	err = ch.QueueBind(
		q.Name,
		util.OrderCreatedTopic,
		util.OrderEventsExchange,
		false,
		nil,
	)
	util.FailOnError(err, fmt.Sprintf("Failed to bind queue: %s to exchange: %s", q.Name, util.OrderEventsExchange))

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, // Make auto ack false and manually acknowledge
		true,
		false,
		false,
		nil,
	)
	util.FailOnError(err, fmt.Sprintf("Failed to consume messages from queue: %s", q.Name))
	for d := range msgs {
		log.Printf("%s", d.Body)
	}
}
