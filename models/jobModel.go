package models

/**
Job体
*/
type Job struct {
	Topic string `json:"topic"`
	ID    string `json:"id"`
	Delay int64  `json:"delay"`
	TTR   int64  `json:"ttr"`
}

/**
关闭vip订单
*/
type CloseOrder struct {
	Job
	Body CloseOrderBody `json:"body"`
}

/**
关闭vip订单具体操作
*/
type CloseOrderBody struct {
	OrderId string `json:"order_id"`
}
