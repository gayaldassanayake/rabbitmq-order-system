# 🧩 Distributed Order Processing System — RabbitMQ Hands-On Project

A hands-on project to master RabbitMQ by building a distributed, event-driven **Order Processing System**.  
The system simulates an e-commerce backend where multiple services communicate asynchronously via RabbitMQ.

---

## 🚀 Overview

When a user places an order, multiple microservices handle their own part of the workflow using **RabbitMQ exchanges, queues, and routing**.

Each service runs independently and communicates through RabbitMQ — no direct REST calls.

### Message Flow

1. Order System
- Receives user’s order request (via HTTP).
- Publishes order.created → order.exchange.

2. Inventory Service
- Consumes order.created.
- Reserves stock, then publishes inventory.reserved or inventory.out_of_stock.

3. Payment Service
- Consumes inventory.reserved.
- Processes payment, then publishes payment.success or payment.failed.

5. Inventory Service (again)
- Consumes payment.success.
- Updates stock levels, then publishes inventory.committed.

6. Order System (again)
- Consumes payment.success.
- Updates order status, then publishes order.paid.

4. Shipping Service
- Consumes payment.success.
- Prepares shipment, then publishes shipping.prepared.

5. Order System (again)
- Consumes shipping.prepared.
- Updates order status, then publishes order.shipped.

6. Shipping Service (again)
- Receives user's confirmation (simulated).
- Publishes shipping.delivered.

7. Order System (again)
- Consumes shipping.delivered.
- Updates order status to order.completed.

6. Audit Service (Consumes from audit.fanout.exchange)
- Consumes all events.
- Logs all events for compliance.

9. Notification Service
- Consumes various order status events (order.*).
- Sends different notifications (email/SMS) using a fanout exchange.

---

## 🐇 RabbitMQ Concepts Covered

| Concept | Used For | Example |
|----------|-----------|---------|
| **Exchanges** (topic, direct, fanout) | Routing messages | `order.exchange`, `notification.exchange` |
| **Queues** | Holding messages | `inventory.queue`, `payment.queue`, etc. |
| **Bindings + Routing Keys** | Selective delivery | `order.created`, `payment.success` |
| **Durable Queues & Persistent Messages** | Reliability | Survive broker restarts |
| **Manual Acknowledgements** | Safe message consumption | `channel.Ack()` |
| **Prefetch (QoS)** | Load control | `channel.Qos(1, false, false)` |
| **Dead-Letter Exchanges (DLX)** | Retry or failure queues | `x-dead-letter-exchange` |
| **Priority Queues** | Urgent messages | `x-max-priority` |
| **Message TTL** | Auto-expire old events | `x-message-ttl` |
| **Publisher Confirms** | Ensure message delivery | `channel.Confirm()` |
| **RPC Pattern** | Request/response (e.g., payment confirmation) | `reply_to`, `correlation_id` |
| **Fanout Exchange** | Broadcast notifications | `notification.exchange` |

---

## 🗂️ Folder Structure

``` 
rabbitmq-order-system/
├── docker-compose.yml # RabbitMQ + services
├── order-service/
│ └── main.go
├── inventory-service/
│ └── main.go
├── payment-service/
│ └── main.go
├── shipping-service/
│ └── main.go
├── notification-service/
│ └── main.go
├── audit-service/
│ └── main.go
└── shared/
├── amqp/
│ ├── connection.go
│ ├── publisher.go
│ └── consumer.go
└── models/
├── order.go
└── event.go
```

## ⚙️ Setup Instructions

### 1️⃣ Prerequisites
- Docker + Docker Compose  
- Go 1.22+  
- RabbitMQ Management UI enabled (port `15672`)

### 2️⃣ Start RabbitMQ
```bash
docker compose up -d
Visit http://localhost:15672
Default login: guest / guest
```

### 3️⃣ Run Each Service

In separate terminals:

``` bash
cd order-service && go run main.go
cd inventory-service && go run main.go
cd payment-service && go run main.go
cd shipping-service && go run main.go
cd notification-service && go run main.go
cd audit-service && go run main.go
```

### 4️⃣ Place an Order

``` bash
curl -X POST http://localhost:8080/orders \
     -H "Content-Type: application/json" \
     -d '{"orderId": "123", "user": "gayal", "items": ["item1", "item2"]}'
```

## 🧪 Example Output

``` bash
[OrderService] Published order.created (OrderID=123)
[InventoryService] Reserved stock for order 123
[PaymentService] Payment success for order 123
[ShippingService] Shipment prepared for order 123
[NotificationService] Sent confirmation email
[AuditService] Logged event: payment.success
```

## 🧭 Suggested Milestones

| Milestone | Focus              | Outcome                                      |
|------------|--------------------|----------------------------------------------|
| 1️⃣ Basic Publish/Consume | Queue, exchange       | Send + receive messages                     |
| 2️⃣ Routing              | Direct/topic          | Selective event delivery                    |
| 3️⃣ Fanout               | Notifications         | Broadcast system                            |
| 4️⃣ Durability           | Persistent queues     | Survive restarts                            |
| 5️⃣ Manual Ack           | Reliability           | Prevent data loss                           |
| 6️⃣ DLX                  | Error handling        | Retry or move to dead-letter queue          |
| 7️⃣ RPC                  | Payment confirmation  | Request–response                            |
| 8️⃣ Monitoring           | Management API        | Track metrics                               |
| 9️⃣ Scaling              | Multiple consumers    | Parallelism and load balancing              |


## 📊 Monitoring & Debugging

Use RabbitMQ Management UI:

- Exchanges → check bindings

- Queues → message rates, unacked counts

- Connections → active services

- Channels → consumer prefetch, acks

For metrics:
```
docker exec -it rabbitmq rabbitmqctl list_queues
```

## 🧠 Learning Outcomes

By completing this project, you will:

- Understand message routing patterns (fanout, topic, direct)

- Implement asynchronous workflows

- Handle failures, retries, and DLQs

- Design loosely coupled microservices

- Gain production-grade RabbitMQ experience

## 🧩 Next Steps (Advanced Ideas)

- Add delayed message retry with plugin

- Integrate Prometheus + Grafana for monitoring

- Add OpenTelemetry tracing

- Containerize all services fully

- Experiment with alternate exchanges and headers exchanges

- Write load tests to observe throughput and bottlenecks

## 🧑‍💻 Author

**Gayal Dassanayake**

Learning RabbitMQ through deep hands-on exploration.
