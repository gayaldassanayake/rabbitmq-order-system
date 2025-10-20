package inventoryservice

import (
	"encoding/json"
	"log"

	"github.com/gayaldassanayake/rabbitmq-order-system/internal/util"
)

func RunService() {
	log.Printf("Inventory service is up and running")
	conn, ch, confirms := util.DeclareRabbitMQChannel()
	defer conn.Close()
	defer ch.Close()

	instockOrders := make(chan util.Order)
	pendingConfirms := make(map[uint64]util.Order)
	go util.VerifyConfirms(confirms, pendingConfirms, instockOrders)

	util.DeclareDomainExchange(ch, util.OrderExchange)
	orders := util.DeclareBindAndConsumeFromQueue(ch, util.OrderCreatedTopic, util.OrderExchange, false)

	util.DeclareDomainExchange(ch, util.InventoryExchange)

	// Start publishing in a goroutine so it can read from instockOrders channel
	go util.PublishEventsFromChannel(
		ch,
		util.InventoryExchange,
		util.InventoryInstockTopic,
		instockOrders,
		pendingConfirms,
	)

	for d := range orders {
		var order util.Order
		err := json.Unmarshal(d.Body, &order)
		if err != nil {
			log.Printf("Error unmarshalling order: %v", err)
			continue
		}
		util.LogStruct(order)
		// TODO: use mongodb and fetch the availability and reserve.
		// For the sake of simplicity we will consider this as always in-stock.
		d.Ack(false)
		instockOrders <- order
	}
}
