package middlewares

import (
	valid "github.com/asaskevich/govalidator"
)

func init() {
	//验证
	valid.SetFieldsRequiredByDefault(true)
}

//保存评价验证
type Remark struct {
	IsCry           int     `json:"is_cry" valid:"-"`
	PracticalRating float64 `json:"practical_rating" valid:"range(2|10)~实用性评分最少为2分，最多为10分"`
	PopularRating   float64 `json:"popular_rating" valid:"range(2|10)~通用性评分最少为2分，最多为10分"`
	LogicRating     float64 `json:"logic_rating" valid:"range(2|10)~逻辑性评分最少为2分，最多为10分"`
	Review          string  `json:"review" valid:"stringlength(1|300)~评论至少1个字，最多300个字"`
}

func (remark *Remark) RemarkValidator() (bool, error) {
	result, err := valid.ValidateStruct(remark)

	return result, err
}
