package util

const (
	OrderServiceUrl  = "/api/v1/order"
	OrderServicePort = ":8080"
	RabbitMQURL      = "amqp://guest:guest@localhost:5672/"
	GlobalRetryCount = 3

	OrderExchange = "order.exchange"
	InventoryExchange = "inventory.exchange"

	OrderCreatedTopic = "order.created.topic"
	InventoryInstockTopic = "inventory.instock.topic"
)
