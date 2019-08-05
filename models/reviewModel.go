package models

/**
按我的风格，评价应该单独建一张表，但是这里却和用户课程放在一张表里面
 */
type Review struct {
	ID              int     `json:"id"`
	Anonymous       int     `json:"anonymous"`
	Rating          float64 `json:"rating"`
	PracticalRating float64 `json:"practical_rating"`
	PopularRating   float64 `json:"popular_rating"`
	LogicRating     float64 `json:"logic_rating"`
	Status          int     `json:"status"`
	Review          string  `json:"review"`
	Reply           string  `json:"reply,omitempty"`
	ReviewedAt      string  `json:"reviewed_at"`
	ReplyAt         string  `json:"reply_at,omitempty"`
	CourseId        int     `json:"course_id"`
}
