package models

/**
购物车
*/
type Cart struct {
	ID       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	CourseId int64  `json:"course_id"`
	Course   Course `json:"course, omitempty"`
}
