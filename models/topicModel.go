package models

/**
题目模型
*/
type TopicModel struct {
	Id            int64         `json:"id"`
	Title         string        `json:"title,omitempty"`
	Type          string        `json:"type,omitempty"`
	Options       string        `json:"options,omitempty"`
	Explan        string        `json:"explan,omitempty"`
	Score         int           `json:"score,omitempty"`
	Difficulty    string        `json:"difficulty,omitempty"`
	Status        string        `json:"status,omitempty"`
	ExtendType    string        `json:"extend_type,omitempty"`
	ExtendContent string        `json:"extend_content,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	UpdatedAt     string        `json:"updated_at,omitempty"`
	CourseId      int64         `json:"course_id,omitempty"`
	ChapterId     int64         `json:"chapter_id,omitempty"`
	ParseOptions  []OptionModel `json:"parse_options,omitempty"` //对Options的处理
}

/**
题目选项
*/
type OptionModel struct {
	Num     string `json:"num"`
	Type    string `json:"type"`
	Content string `json:"content"`
	IsRight string `json:"is_right"`
}
