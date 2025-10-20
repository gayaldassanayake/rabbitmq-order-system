package orderservice

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gayaldassanayake/rabbitmq-order-system/internal/util"
)

func RunService() {
	log.Printf("Order service is up and running")
	orderChan := make(chan util.Order)
	// TODO: persist order in MongoDB
	go publishOrderEvents(orderChan)

	http.HandleFunc(util.OrderServiceUrl, orderHandler(orderChan))
	http.ListenAndServe(util.OrderServicePort, nil)
}

func orderHandler(ch chan<- util.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		}

		var orderReq util.OrderRequest
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		order := util.Order{
			Id:           util.GenerateUUID(),
			OrderRequest: &orderReq,
		}
		order.Id = util.GenerateUUID()
		ch <- order

		w.WriteHeader(http.StatusCreated)
	}
}

func publishOrderEvents(httpCh chan util.Order) {
	conn, ch, confirms := util.DeclareRabbitMQChannel()
	defer conn.Close()
	defer ch.Close()

	pendingConfirms := make(map[uint64]util.Order)
	go util.VerifyConfirms(confirms, pendingConfirms, httpCh)

	util.DeclareDomainExchange(ch, util.OrderExchange);

	util.PublishEventsFromChannel(ch, util.OrderExchange, util.OrderCreatedTopic, httpCh, pendingConfirms)
}
