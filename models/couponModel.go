package models

import "time"

/**
优惠券
*/
type Coupon struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Value         float64   `json:"value"`
	Total         int       `json:"total"`
	LimitNumber   int       `json:"limit_number"`
	UsedNumber    int       `json:"used_number"`
	Suitable      string    `json:"suitable"`
	SuitableValue int       `json:"suitable_value"`
	MinAmount     float64   `json:"min_amount"`
	NotBefore     time.Time `json:"not_before"`
	NotAfter      time.Time `json:"not_after"`
	EffectiveDay  int       `json:"effective_day"`
	Enabled       int       `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
