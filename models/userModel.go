package models

type User struct {
	Id               int         `json:"id"`
	Avatar           string      `json:"avatar"`
	No               string      `json:"no,omitempty"`
	Nickname         string      `json:"nickname"`
	Title            interface{} `json:"title,omitempty"`
	Mobile           string      `json:"mobile,omitempty"`
	Email            string      `json:"email,omitempty"`
	Password         string      `json:"password,omitempty"`
	MobileVerified   int         `json:"mobile_verified,omitempty"`
	MobileVerifiedAt string      `json:"mobile_verified_at,omitempty"`
	EmailVerified    int         `json:"email_verified,omitempty"`
	EmailVerifiedAt  string      `json:"email_verified_at,omitempty"`
	WechatVerified   int         `json:"wechat_verified,omitempty"`
	WechatVerifiedAt string      `json:"wechat_verified_at,omitempty"`
	Level            string      `json:"level,omitempty"`
	StartAt          string      `json:"start_at,omitempty"`
	EndAt            string      `json:"end_at,omitempty"`
	Status           int         `json:"status,omitempty"`
	IsLecturer       int         `json:"is_lecturer,omitempty"`
	About            string      `json:"about,omitempty"`
	Source           string      `json:"source,omitempty"`
	RegisterAt       string      `json:"register_at,omitempty"`
	RegisterIp       string      `json:"register_ip,omitempty"`
	RegisterCity     string      `json:"register_city,omitempty"`
	LastLoginAt      string      `json:"last_login_at,omitempty"`
	LastLoginIp      string      `json:"last_login_ip,omitempty"`
	LastLoginCity    string      `json:"last_login_city,omitempty"`
	LoginCount       int         `json:"login_count,omitempty"`
	CreatedAt        string      `json:"created_at,omitempty"`
	UpdatedAt        string      `json:"updated_at,omitempty"`
}
