package models

/**
题目模型
*/
type TopicModel struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	Options       string `json:"options,omitempty"`
	Explan        string `json:"explan,omitempty"`
	Score         int    `json:"score"`
	Difficulty    string `json:"difficulty"`
	Status        string `json:"status"`
	ExtendType    string `json:"extend_type,omitempty"`
	ExtendContent string `json:"extend_content,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CourseId      int64  `json:"course_id"`
	ChapterId     int64  `json:"chapter_id,omitempty"`
}
