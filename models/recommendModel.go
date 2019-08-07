package models

/**
推荐课程
 */
type Recommend struct {
	Id            int     `json:"id"`
	Type          string  `json:"type"`
	Title         string  `json:"title"`
	Price         float32 `json:"price"`
	VipPrice      float32 `json:"vip_price"`
	Discount      float32 `json:"discount"`
	DiscountEndAt string  `json:"discount_end_at"`
	CoverPicture  string  `json:"cover_picture"`
	LearnNum      int     `json:"learn_num"`
	BuyNum        int     `json:"buy_num"`
}
