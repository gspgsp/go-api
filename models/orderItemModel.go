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
	UpdatedAt    string  `json:"updated_at"`
	OrderId      int     `json:"order_id"`
	PeriodId     int     `json:"period_id"`
	TrainingId   int     `json:"training_id"`
	CourseId     int     `json:"course_id"`
	UserId       int     `json:"user_id"`
}
