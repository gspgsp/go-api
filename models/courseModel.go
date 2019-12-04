package models

type Course struct {
	Id              int     `json:"id"`
	Type            string  `json:"type,omitempty"`
	Title           string  `json:"title,omitempty"`
	SubTitle        string  `json:"subtitle,omitempty"`
	Price           float32 `json:"price,omitempty"`
	VipPrice        float32 `json:"vip_price,omitempty"`
	Discount        float32 `json:"discount,omitempty"`
	DiscountEndAt   string  `json:"discount_end_at,omitempty"`
	CoverPicture    string  `json:"cover_picture,omitempty"`
	BackPicture     string  `json:"back_picture,omitempty"`
	LearnNum        int     `json:"learn_num,omitempty"`
	BuyNum          int     `json:"buy_num,omitempty"`
	VideoUrl        string  `json:"video_url,omitempty"`
	FavNum          int     `json:"fav_num,omitempty"`
	VipLevel        int     `json:"vip_level,omitempty"`
	DifficultyLevel int     `json:"difficulty_level,omitempty"`
	IsRecommended   int     `json:"is_recommended,omitempty"`
	Length          int     `json:"length,omitempty"`
	LessonNum       int     `json:"lesson_num,omitempty"`
	EffectiveDay    int     `json:"effective_day,omitempty"`
	ServiceQrCode   string  `json:"service_qr_code,omitempty"`
	Rating          float32 `json:"rating,omitempty"`
	PracticalRating float32 `json:"practical_rating,omitempty"`
	PopularRating   float32 `json:"popular_rating,omitempty"`
	LogicRating     float32 `json:"logic_rating,omitempty"`
	ReviewCount     int     `json:"review_count,omitempty"`
	Keywords        string  `json:"keywords,omitempty"`
	Description     string  `json:"description,omitempty"`
	Goals           string  `json:"goals,omitempty"`
	Audiences       string  `json:"audiences,omitempty"`
	Summary         string  `json:"summary,omitempty"`
	PcBack          string  `json:"pc_back,omitempty"`
	PcAbout         string  `json:"pc_about,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
	CategoryId      int     `json:"category_id,omitempty"`
}

func (Course) TableName() string {
	return "h_edu_courses"
}

/**
拓展信息
*/
type OrderCourse struct {
	PeriodInfo   *Period   `json:"period_info,omitempty"`
	TrainingInfo *Training `json:"training_info,omitempty"`
	PackageInfo  *Package  `json:"package_info,omitempty"`
	Courses      []Course  `json:"courses"`
}
