package models

type UserCourse struct {
	Id              int      `json:"id"`
	Type            string   `json:"type"`
	CreatedAt       string   `json:"created_at"`
	FinishedAt      string   `json:"finished_at"`
	Reviewed        int      `json:"reviewed"`
	Anonymous       int      `json:"anonymous"`
	Rating          float32  `json:"rating"`
	PracticalRating int      `json:"practical_rating"`
	PopularRating   int      `json:"popular_rating"`
	LogicRating     int      `json:"logic_rating"`
	Status          int      `json:"status"`
	Review          string   `json:"review"`
	ReviewedAt      string   `json:"reviewed_at"`
	Schedule        int      `json:"schedule"`
	CourseId        int      `json:"course_id"`
	LessonId        int      `json:"lesson_id"`
	UserId          int      `json:"user_id"`
	Course          []Course `json:"course"`
}

func (UserCourse) TableName() string {
	return "h_user_course"
}
