package models

/**
公告模型
*/
type NoticeModel struct {
	ID             int64  `json:"id"`
	Type           int    `json:"type"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	StartAt        int64  `json:"start_at"`
	EndAt          int64  `json:"end_at"`
	Status         int    `json:"status"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	CreatedAdminId int    `json:"created_admin_id, omitempty"`
}
