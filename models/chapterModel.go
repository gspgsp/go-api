package models

type Chapter struct {
	Id          int64     `json:"id"`
	ParentId    int64     `json:"parent_id"`
	Type        string    `json:"type"`
	Number      int64     `json:"number,omitempty"`
	NumberPath  string    `json:"number_path,omitempty"`
	Title       string    `json:"title,omitempty"`
	Preview     string    `json:"preview"`
	IsFree      int       `json:"is_free"`
	LessonType  string    `json:"lesson_type"`
	MediaSource string    `json:"media_source"`
	Path        string    `json:"path,omitempty"`
	Length      int       `json:"length,omitempty"`
	Status      string    `json:"status"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at,omitempty"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
	CourseId    int64     `json:"course_id"`
	Children    []Chapter `json:"children,omitempty"` //子类
}
