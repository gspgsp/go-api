package models

type Course struct {
	Id              int     `json:"id"`
	Type            string  `json:"type"`
	Title           string  `json:"title"`
	SubTitle        string  `json:"subtitle"`
	Price           float32 `json:"price"`
	VipPrice        float32 `json:"vip_price"`
	Discount        float32 `json:"discount"`
	DiscountEndAt   string  `json:"discount_end_at"`
	CoverPicture    string  `json:"cover_picture"`
	BackPicture     string  `json:"back_picture"`
	LearnNum        int     `json:"learn_num"`
	BuyNum          int     `json:"buy_num"`
	VideoUrl        string  `json:"video_url"`
	FavNum          int     `json:"fav_num"`
	VipLevel        int     `json:"vip_level"`
	DifficultyLevel int     `json:"difficulty_level"`
	IsRecommended   int     `json:"is_recommended"`
	Length          int     `json:"length"`
	LessonNum       int     `json:"lesson_num"`
	EffectiveDay    int     `json:"effective_day"`
	ServiceQrCode   string  `json:"service_qr_code"`
	Rating          float32 `json:"rating"`
	PracticalRating float32 `json:"practical_rating"`
	PopularRating   float32 `json:"popular_rating"`
	LogicRating     float32 `json:"logic_rating"`
	ReviewCount     int     `json:"review_count"`
	Keywords        string  `json:"keywords"`
	Description     string  `json:"description"`
	Goals           string  `json:"goals"`
	Audiences       string  `json:"audiences"`
	Summary         string  `json:"summary"`
	PcBack          string  `json:"pc_back"`
	PcAbout         string  `json:"pc_about"`
	CreatedAt       string  `json:"created_at"`
	CategoryId      int     `json:"category_id"`
}

func (Course) TableName() string  {
	return "h_edu_course"
}
