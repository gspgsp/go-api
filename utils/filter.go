package utils

/**
接口参数对象
 */
type Filter struct {
	Type            string `json:"type"`
	Code            string `json:"code"`
	DifficultyLevel string `json:"difficulty_level"`
	Order           string `json:"order"`
	VipPrice        string `json:"vip_price"`
	IsRecommended   string `json:"is_recommended"`
}

/**
真实IP信息信息
 */
type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data"`
}

type IP struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Isp       string `json:"isp"`
}