package models

type User struct {
	Id               int    `json:"id"`
	Avatar           string `json:"avatar"`
	No               string `json:"no"`
	Nickname         string `json:"nickname"`
	Title            int    `json:"title"`
	Mobile           string `json:"mobile"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	MobileVerified   int    `json:"mobile_verified"`
	MobileVerifiedAt string `json:"mobile_verified_at"`
	EmailVerified    int    `json:"email_verified"`
	EmailVerifiedAt  string `json:"email_verified_at"`
	WechatVerified   int    `json:"wechat_verified"`
	WechatVerifiedAt string `json:"wechat_verified_at"`
	Level            string `json:"level"`
	StartAt          string `json:"start_at"`
	EndAt            string `json:"end_at"`
	Status           int    `json:"status"`
	IsLecturer       int    `json:"is_lecturer"`
	About            string `json:"about"`
	Source           string `json:"source"`
	RegisterAt       string `json:"register_at"`
	RegisterIp       string `json:"register_ip"`
	RegisterCity     string `json:"register_city"`
	LastLoginAt      string `json:"last_login_at"`
	LastLoginIp      string `json:"last_login_ip"`
	LastLoginCity    string `json:"last_login_city"`
	LoginCount       int    `json:"login_count"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
