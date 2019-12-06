package models

import "time"

/**
用户优惠券
*/
type UserCoupon struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Suitable      string    `json:"suitable"`
	SuitableValue int       `json:"suitable_value"`
	Status        int       `json:"status"`
	UsedAt        string    `json:"used_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserId        int       `json:"user_id"`
	CouponId      int       `json:"coupon_id"`
}

/**
优惠券处理
*/
type CouponInfo struct {
	UserCoupon
	CName         string    `json:"c_name"`
	CValue        float32   `json:"c_value"`
	CMinAmount    float32   `json:"c_min_amount"`
	CSuitable     string    `json:"c_suitable"`
	CNotBefore    time.Time `json:"c_not_before"`
	CNotAfter     time.Time `json:"c_not_after"`
	CEffectiveDay int       `json:"c_effective_day"`
	CEnabled      int       `json:"c_enabled"`
}
