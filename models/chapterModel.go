package models

type Chapter struct {
	Id          int64  `json:"id"`
	ParentId    int64  `json:"parent_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Preview     string `json:"preview"`
	IsFree      int    `json:"is_free"`
	LessonType  string `json:"lesson_type"`
	MediaSource string `json:"media_source"`
	Length      int    `json:"length"`
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
