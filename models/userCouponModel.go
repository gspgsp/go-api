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
	UsedAt        time.Time `json:"used_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserId        int       `json:"user_id"`
	CouponId      int       `json:"coupon_id"`
	CouponInfo    Coupon    `json:"coupon_info,omitempty"`
}
