package models

import "time"

/**
训练营模型
*/
type Training struct {
	ID           int       `json:"id"`
	Uuid         string    `json:"uuid"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	ShowContent  string    `json:"show_content"`
	CoverPicture string    `json:"cover_picture"`
	LearnNum     int       `json:"learn_num"`
	BuyNum       int       `json:"buy_num"`
	Status       string    `json:"status"`
	PcBack       string    `json:"pc_back"`
	PcAbout      string    `json:"pc_about"`
	MbAbout      string    `json:"mb_about"`
	BuyNotice    string    `json:"buy_notice"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
