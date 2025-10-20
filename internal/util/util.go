package util

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func LogStruct(v interface{}) {
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

func VerifyConfirms[T any](confirms <-chan amqp.Confirmation, pendingConfirms map[uint64]T, retryChan chan<- T) {
	for confirm := range confirms {
		if item, exists := pendingConfirms[confirm.DeliveryTag]; exists {
			if !confirm.Ack {
				go onConfirmFailure(item, retryChan)
			}
		}
		delete(pendingConfirms, confirm.DeliveryTag)
	}
}

func GetMockProbResponse(prob int) bool {
	return rand.Intn(100) < prob
}

func DeclareDomainExchange(ch *amqp.Channel, name string) {
	err := ch.ExchangeDeclare(
		name,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, fmt.Sprintf("Failed to declare exchange: %s", name))
}

func DeclareBindAndConsumeFromQueue(ch *amqp.Channel, topicName, exchangeName string, autoack bool) <-chan amqp.Delivery {
	q, err := ch.QueueDeclare(
		"",
		true,
		false,
		true,
		false,
		nil,
	)
	FailOnError(err, fmt.Sprintf("Failed to declare queue: %s", q.Name))

	err = ch.QueueBind(
		q.Name,
		topicName,
		exchangeName,
		false,
		nil,
	)
	FailOnError(err, fmt.Sprintf("Failed to bind queue: %s to exchange: %s", q.Name, exchangeName))

	msgs, err := ch.Consume(
		q.Name,
		"",
		autoack,
		true,
		false,
		false,
		nil,
	)
	FailOnError(err, fmt.Sprintf("Failed to consume messages from queue: %s", q.Name))
	return msgs
}

func PublishEventsFromChannel[T any](
	ch *amqp.Channel, 
	exchangeName string, 
	topicName string, 
	inputChan <-chan T,
	pendingConfirms map[uint64]T,
	) {
	for t := range inputChan {
		body, err := json.Marshal(t)
		if err != nil {
			log.Printf("Failed to marshal order: %v", err)
			continue
		}
		DeliveryTag := ch.GetNextPublishSeqNo()
		err = ch.Publish(
			exchangeName,
			topicName,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: GenerateUUID(),
				DeliveryMode:  2,
				Body:          body,
			},
		)
		if err == nil {
			pendingConfirms[DeliveryTag] = t
		} else {
			// TODO: add this to the retry mechanism
		}
	}
}

func onConfirmFailure[T any](retryItem T, retryChan chan<- T) {
	time.Sleep(500 * time.Millisecond)
	select {
	case retryChan <- retryItem:
	default:
		log.Printf("Failed to re-queue order (channel full): %v", retryItem)
	}
}
