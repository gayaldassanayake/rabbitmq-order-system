package util

const (
	OrderServiceUrl  = "/api/v1/order"
	OrderServicePort = ":8080"
	RabbitMQURL      = "amqp://guest:guest@localhost:5672/"
	GlobalRetryCount = 3

	OrderEventsExchange = "order.events.exchange"

	OrderCreatedTopic = "order.created.topic"
)
