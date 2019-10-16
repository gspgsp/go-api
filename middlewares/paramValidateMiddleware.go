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

//保存答题验证
type Answer struct {
	CourseId  int64        `json:"course_id" valid:"required~课程id必须"`
	RollId    int64        `json:"roll_id" valid:"required~试卷id必须"`
	StartTime int64        `json:"start_time" valid:"required~答题用时必须"`
	Answers   []AnswerData `json:"answer_data" valid:"required~答案必须"`
}

//具体答案选项
type AnswerData struct {
	TopicId int64  `json:"topic_id" valid:"required~题目id必须"`
	Option  string `json:"option" valid:"required~当前题目选项必须"`
}

func (answer *Answer) AnswerValidator() (bool, error) {
	//自定义验证规则，本来想自定义一个针对结构体嵌套的验证的，但是发现木有用；其实可以直接定义验证规则，验证器会自动嵌套验证，如上[]AnswerData
	//valid.CustomTypeTagMap.Set("answerDataValidator", func(i, context interface{}) bool {
	//	switch v := i.(type) {
	//	case []AnswerData:
	//		for _, e := range v {
	//			res, _ := valid.ValidateStruct(e)
	//			return res
	//		}
	//	}
	//	return false
	//})

	result, err := valid.ValidateStruct(answer)
	return result, err
}

//创建VIP订单验证
type VipOrder struct {
	Id     int64  `json:"id" valid:"in(1)~vip id 必须为1"`
	Source string `json:"source" valid:"-"`
}

func (vipOrder *VipOrder) VipOrderValidator() (bool, error) {
	result, err := valid.ValidateStruct(vipOrder)
	return result, err
}
