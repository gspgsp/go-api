package models

type SlideModel struct {
	ID             int64  `json:"id"`
	Port           int    `json:"port"`
	Title          string `json:"title"`
	Url            string `json:"url"`
	Carousel       string `json:"carousel"`
	Sort           int    `json:"sort,omitempty"`
	Status         int    `json:"status"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	CreatedAdminId int    `json:"created_admin_id,omitempty"`
}
