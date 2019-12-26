package models

/**
支付宝支付拓展参数
*/
type PayAliExtendParam struct {
	BranchType string `json:"branch_type"`
	Id         int    `json:"id"`
	PaySource  string `json:"pay_source"`
}
