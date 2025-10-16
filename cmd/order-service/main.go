package main

import (
	orderservice "github.com/gayaldassanayake/rabbitmq-order-system/order-service"
)

func main() {
	orderservice.RunService()
}
