package orderservice

type OrderRequest struct {
	UserId	string	`json:"userId"`
	Items	[]Item	`json:"items"`
}

type Order struct {
	*OrderRequest
	Id string
}

type Item struct {
	Code	string	`json:"code"`
	UnitPrice	float64	`json:"unit_price"`
	Quantity	float64	`json:"quantity"`
}
