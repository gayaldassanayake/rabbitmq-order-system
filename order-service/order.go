package orderservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gayaldassanayake/rabbitmq-order-system/internal/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RunService() {
	log.Printf("Order service is up and running")
	orderChan := make(chan Order)
	// TODO: persist order in the DB
	go publishOrderEvents(orderChan)

	http.HandleFunc(util.OrderServiceUrl, orderHandler(orderChan))
	http.ListenAndServe(util.OrderServicePort, nil)
}

func orderHandler(ch chan<- Order) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		}

		var orderReq OrderRequest
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		order := Order{
			Id: util.GenerateUUID(),
			OrderRequest: &orderReq,
		}
		order.Id = util.GenerateUUID()
		ch <- order

		w.WriteHeader(http.StatusCreated)
	}
}

func publishOrderEvents(httpCh chan Order) {
	conn, ch, confirms := util.DeclareRabbitMQChannel()
	defer conn.Close()
	defer ch.Close()
	
	pendingConfirms := make(map[uint64]Order)
	go func() {
		for confirm := range confirms {
			if order, exists := pendingConfirms[confirm.DeliveryTag]; exists {
				if !confirm.Ack {
					go func(retryOrder Order) {
						time.Sleep(500 * time.Millisecond)
						select {
						case httpCh <- retryOrder:
							log.Printf("Order re-queued for retry :%s", retryOrder.Id)
						default:
							log.Printf("Failed to re-queue order (channel full): %s", retryOrder.Id)
						}
					}(order)
				}
			}
			delete(pendingConfirms, confirm.DeliveryTag)
		}
	}()

	// declare exchange, queue
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

	for order := range httpCh {
		body, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order: %v", err)
			continue
		}
		DeliveryTag := ch.GetNextPublishSeqNo()
		err = ch.Publish(
			util.OrderEventsExchange,
			util.OrderCreatedTopic,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: util.GenerateUUID(),
				DeliveryMode:  2,
				Body:          body,
			},
		)
		if err == nil {
			pendingConfirms[DeliveryTag] = order
		} else {
			// TODO: add this to the retry mechanism somehow
		}
	}
}
