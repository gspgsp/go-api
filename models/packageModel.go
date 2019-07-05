package models

type Package struct {
	Id            int     `json:"id"`
	Type          string  `json:"type"`
	Title         string  `json:"title"`
	Subtitle      string  `json:"subtitle"`
	Price         float32 `json:"price"`
	Discount      float32 `json:"discount"`
	DiscountEndAt string  `json:"discount_end_at"`
	CoverPicture  string  `json:"cover_picture"`
	BackPicture   string  `json:"back_picture"`
	LearnNum      int     `json:"learn_num"`
	BuyNum        int     `json:"buy_num"`
	VipLevel      int     `json:"vip_level"`
	VideoUrl      string  `json:"video_url"`
	Keywords      string  `json:"keywords"`
	Description   string  `json:"description"`
	Goals         []byte  `json:"goals"`
	Audiences     []byte  `json:"audiences"`
	MbAbout       []byte  `json:"mb_about"`
	CreatedAt     string  `json:"created_at"`
}
