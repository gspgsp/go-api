package models

import "time"

type VipModel struct {
	ID             int64       `json:"id"`
	Title          string      `json:"title"`
	Subtitle       string      `json:"subtitle,omitempty"`
	Price          float64     `json:"price"`
	Discount       float64     `json:"discount,omitempty"`
	DiscountEndAt  string      `json:"discount_end_at,omitempty"`
	LearnNum       int64       `json:"learn_num,omitempty"`
	BuyNum         int64       `json:"buy_num,omitempty"`
	VideoUrl       string      `json:"video_url,omitempty"`
	VipLevel       string      `json:"vip_level,omitempty"`
	EffectiveDay   int64       `json:"effective_day,omitempty"`
	ServiceQrCode  string      `json:"service_qr_code,omitempty"`
	Keywords       string      `json:"keywords,omitempty"`
	Description    string      `json:"description,omitempty"`
	PcAbout        string      `json:"pc_about,omitempty"`
	MbAbout        string      `json:"mb_about,omitempty"`
	Problem        string      `json:"problem,omitempty"`
	CreatedAt      time.Time   `json:"created_at,omitempty"`
	ParseCreatedAt interface{} `json:"parse_created_at"`
	UpdatedAt      time.Time   `json:"updated_at,omitempty"`
	ParseUpdatedAt interface{} `json:"parse_updated_at"`
	IsBuy          int         `json:"is_buy"` //判断当前用户是否买过这个会员
}
