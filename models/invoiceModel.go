package models

type InvoiceModel struct {
	ID               int     `json:"id"`
	No               string  `json:"no"`
	Body             string  `json:"body"`
	InvoiceType      string  `json:"invoice_type"`
	InvoiceNo        string  `json:"invoice_no"`
	InvoiceAmount    float32 `json:"invoice_amount"`
	InvoiceStatus    int     `json:"invoice_status"`
	InvoiceAt        string  `json:"invoice_at"`
	PaymentOrderNo   string  `json:"payment_order_no"`
	PaymentAmount    float32 `json:"payment_amount"`
	ReceiptAmount    float32 `json:"receipt_amount"`
	PaymentMethod    int     `json:"payment_method"`
	PaymentStatus    int     `json:"payment_status"`
	PaymentExpiredAt string  `json:"payment_expired_at"`
	PaymentAt        string  `json:"payment_at"`
	Source           string  `json:"source"`
	Status           int     `json:"status"`
	Company          string  `json:"company"`
	NorCode          string  `json:"nor_code"`
	AddressMobile    string  `json:"address_mobile"`
	BankNumber       string  `json:"bank_number"`
	Content          string  `json:"content"`
	ContactName      string  `json:"contact_name"`
	ContactMobile    string  `json:"contact_mobile"`
	ContactEmail     string  `json:"contact_email"`
	Address          string  `json:"address"`
	ZipCode          string  `json:"zip_code"`
	ExpType          string  `json:"exp_type"`
	ExpName          string  `json:"exp_name"`
	ExpNumber        string  `json:"exp_number"`
	UserRemark       string  `json:"user_remark"`
	AdminRemark      string  `json:"admin_remark"`
	Extra            string  `json:"extra"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	UserId           int     `json:"user_id"`
}
