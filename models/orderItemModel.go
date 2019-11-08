package models

/**
订单详情模型
*/
type OrderItemModel struct {
	ID           int     `json:"id"`
	OType        string  `json:"type"`
	Price        float32 `json:"price"`
	PaymentPrice float32 `json:"payment_price"`
	CreatedAt    string  `json:"created_at"`
	OrderId      int64   `json:"order_id"`
	PeriodId     int64   `json:"period_id"`
	TrainingId   int64   `json:"training_id"`
	CourseId     int64   `json:"course_id"`
	UserId       int64   `json:"user_id"`
}
