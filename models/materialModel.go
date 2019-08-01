package models

type Material struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	FileName    string  `json:"file_name,omitempty"`
	Description string  `json:"description,omitempty"`
	Link        string  `json:"link,omitempty"`
	Size        float64 `json:"size,omitempty"`//置空以后不显示
	FormatSize  string  `json:"format_size,omitempty"`//为了接受格式化以后的size
	Type        string  `json:"type"`
	DownloadNum int     `json:"download_num"`
	CreatedAt   string  `json:"created_at,omitempty"`
	CourseId    int     `json:"course_id"`
}
