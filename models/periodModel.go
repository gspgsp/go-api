package models

import "time"

/**
期模型
*/
type Period struct {
	ID           int       `json:"id"`
	Uuid         string    `json:"uuid"`
	Number       int64     `json:"number"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	SignUpEndAt  time.Time `json:"sign_up_end_at"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	MasterWechat string    `json:"master_wechat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TrainingId   int       `json:"training_id"`
	CourseId     int       `json:"course_id"`
}
