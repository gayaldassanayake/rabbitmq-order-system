package paymentservice

import (
	"encoding/json"
	"log"

	"github.com/gayaldassanayake/rabbitmq-order-system/internal/util"
)

func RunService() {
	log.Printf("Payment service is up and running")
	conn, ch, confirms := util.DeclareRabbitMQChannel()
	defer conn.Close()
	defer ch.Close()

	payedOrders := make(chan util.Order)
	pendingConfirms := make(map[uint64]util.Order)
	go util.VerifyConfirms(confirms, pendingConfirms, payedOrders)

	util.DeclareDomainExchange(ch, util.InventoryExchange)
	orders := util.DeclareBindAndConsumeFromQueue(ch, util.InventoryInstockTopic, util.InventoryExchange, false)

	for d := range orders {
		var order util.Order
		err := json.Unmarshal(d.Body, &order)
		if err != nil {
			log.Printf("Error unmarshalling order: %v", err)
			continue
		}
		// TODO: For the sake of simplicity we will consider payment to be always successful.
		d.Ack(false)
		util.LogStruct(order)
		payedOrders <- order
	}
}
