package models

/**
题库模型
*/
type RollModel struct {
	Id              int64  `json:"id"`
	Type            string `json:"type"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	TotalScore      int64  `json:"total_score,omitempty"`
	Mode            int    `json:"mode"`
	ItemCount       int64  `json:"item_count"`
	PassedCondition string `json:"passed_condition,omitempty"`
	LimitedAt       int64  `json:"limited_at"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	CourseId        int64  `json:"course_id"`
	ChapterId       int64  `json:"chapter_id,omitempty"`
}
