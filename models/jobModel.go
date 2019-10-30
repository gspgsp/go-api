package models

type CloseOrder struct {
	Topic string         `json:"topic"`
	ID    string         `json:"id"`
	Delay string         `json:"delay"`
	TTR   string         `json:"ttr"`
	Body  CloseOrderBody `json:"body"`
}

type CloseOrderBody struct {
	OrderId string `json:"order_id"`
}
