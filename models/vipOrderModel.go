package models

import "time"

type VipOrderModel struct {
	ID                   int       `json:"id"`
	No                   string    `json:"no"`
	Amount               float64   `json:"amount"`
	DiscountAmount       float64   `json:"discount_amount"`
	PaymentOrderNo       string    `json:"payment_order_no"`
	PaymentAmount        float64   `json:"payment_amount"`
	ReceiptAmount        float64   `json:"receipt_amount"`
	PaymentMethod        int       `json:"payment_method"`
	PaymentStatus        int       `json:"payment_status"`
	PaymentExpiredAt     time.Time `json:"payment_expired_at"`
	PaymentAt            time.Time `json:"payment_at"`
	Source               string    `json:"source"`
	Status               int       `json:"status"`
	Reviewed             int       `json:"reviewed,omitempty"`
	RefundReason         string    `json:"refund_reason,omitempty"`
	RefundRequestAt      time.Time `json:"refund_request_at,omitempty"`
	RefundDisagreeReason string    `json:"refund_disagree_reason,omitempty"`
	RefundStatus         string    `json:"refund_status,omitempty"`
	RefundNo             string    `json:"refund_no,omitempty"`
	RefundAt             time.Time `json:"refund_at,omitempty"`
	Extra                string    `json:"extra,omitempty"`
	UserRemark           string    `json:"user_remark,omitempty"`
	AdminRemark          string    `json:"admin_remark,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UserId               int64     `json:"user_id"`
	VipId                int64     `json:"vip_id"`
	InvoiceId            int64     `json:"invoice_id,omitempty"`
}
