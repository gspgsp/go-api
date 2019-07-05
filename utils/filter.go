package utils

type Filter struct {
	Type            string `json:"type"`
	Code            string `json:"code"`
	DifficultyLevel string `json:"difficulty_level"`
	Order           string `json:"order"`
	VipPrice        string `json:"vip_price"`
	IsRecommended   string `json:"is_recommended"`
}
