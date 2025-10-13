package orderservice

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gayaldassanayake/rabbitmq-order-system/internal/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

const urlPrefix = "/api/v1/order"

func RunOrderService() {
	log.Printf("Order service is up and running")
	orderChan := make(chan Order)
	go publishOrderEvents(orderChan)

	http.HandleFunc(urlPrefix, orderHandler(orderChan))
	http.ListenAndServe(":8080", nil)
}

func orderHandler(ch chan <- Order) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		}

		var order Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ch <- order

		w.WriteHeader(http.StatusCreated)
	}
}

func publishOrderEvents(httpCh <- chan Order) {
	for order := range httpCh {
		util.LogStruct(order)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.FailOnError(err, "Failed to create rabbitmq connection")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	
}
