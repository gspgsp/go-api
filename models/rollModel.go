package models

/**
题库模型
*/
type RollModel struct {
	Id              int64      `json:"id"`
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Description     string     `json:"description,omitempty"`
	TotalScore      int        `json:"total_score"`
	Mode            int        `json:"mode"`
	ItemCount       int64      `json:"item_count"`
	PassedCondition string     `json:"passed_condition,omitempty"`
	LimitedAt       int64      `json:"limited_at"`
	Status          string     `json:"status"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
	CourseId        int64      `json:"course_id"`
	ChapterId       int64      `json:"chapter_id,omitempty"`
	Grade           GradeModel `json:"grade,omitempty"`
}

/**
题库题目详情
*/
type RollInfoModel struct {
	Id        int64        `json:"id"`
	Title     string       `json:"title"`
	ItemCount int64        `json:"item_count"`
	LimitedAt int64        `json:"limited_at"`
	CourseId  int64        `json:"course_id"`
	Topics    []TopicModel `json:"topics"`
}
