package models

/**
订单模型
*/
type OrderModel struct {
	ID               int     `json:"id"`
	No               string  `json:"no"`
	Amount           float32 `json:"amount"`
	DiscountAmount   float32 `json:"discount_amount"`
	CouponAmount     float32 `json:"coupon_amount"`
	PaymentOrderNo   string  `json:"payment_order_no"`
	PaymentAmount    float32 `json:"payment_amount"`
	ReceiptAmount    float32 `json:"receipt_amount"`
	PaymentMethod    int     `json:"payment_method"`
	PaymentStatus    int     `json:"payment_status"`
	PaymentExpiredAt string  `json:"payment_expired_at"`
	PaymentAt        string  `json:"payment_at"`
	Source           string  `json:"source"`
	Status           int     `json:"status"`
	RefundReason     string  `json:"refund_reason"`
	RefundRequestAt  string  `json:"refund_request_at"`
	RefundStatus     string  `json:"refund_status"`
	RefundNo         string  `json:"refund_no"`
	RefundAt         string  `json:"refund_at"`
	Extra            string  `json:"extra"`
	UserRemark       string  `json:"user_remark"`
	AdminRemark      string  `json:"admin_remark"`
	CreatedAt        string  `json:"created_at"`
	UserId           int64   `json:"user_id"`
	UserCouponId     int64   `json:"user_coupon_id"`
	InvoiceId        int64   `json:"invoice_id"`
	PackageId        int64   `json:"package_id"`
}
