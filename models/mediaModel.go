package models

type Media struct {
	Chapter       Chapter
	CurrentTitle  string `json:"current_title"`
	CurrentLesion string `json:"current_lesion"`
	LesionType    string `json:"lesion_type"`
}
