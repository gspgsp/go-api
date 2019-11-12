package models

/**
购物车
*/
type Cart struct {
	ID       int    `json:"id"`
	UserId   int    `json:"user_id"`
	CourseId int    `json:"course_id"`
	Course   Course `json:"course, omitempty"`
}
