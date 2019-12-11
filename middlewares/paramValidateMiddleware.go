package middlewares

import (
	"errors"
	valid "github.com/asaskevich/govalidator"
	"strconv"
	"strings"
	"sync"
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

//添加购物车验证
type AddCart struct {
	CourseId int64 `json:"course_id" valid:"required~课程id必须"`
}

func (addCart *AddCart) AddCartValidator() (bool, error) {
	result, err := valid.ValidateStruct(addCart)
	return result, err
}

//提交订单验证
type CommitOrder struct {
	Source       string `json:"source" valid:"-"`
	Type         string `json:"type" valid:"-"`
	Ids          string `json:"ids" valid:"legalId"`
	PeriodId     int    `json:"period_id" valid:"-"`
	TrainingId   int    `json:"training_id" valid:"-"`
	UserCouponId int    `json:"user_coupon_id" valid:"-"`
	ChannelUuid  string `json:"channel_uuid" valid:"-"`
	UserMark     string `json:"user_mark" valid:"-"`
}

func (commitOrder *CommitOrder) CommitOrderValidator() (bool, error) {
	//本来想用struct的in查询的，发现in(字符串不行)
	if !valid.IsIn(commitOrder.Source, "pc", "mb") || !valid.IsIn(commitOrder.Type, "course", "package", "training") {
		return false, errors.New("客户端来源/课程类型不正确")
	}

	if commitOrder.Type == "training" && (commitOrder.PeriodId == 0 && commitOrder.TrainingId == 0) {
		return false, errors.New("训练营/期ID必须")
	}

	var mt sync.Mutex
	mt.Lock()
	//这里如果直接用 services.BaseOrm的话会报错误: import cycle not allowed，这里直接验证ids是否合理
	valid.TagMap["legalId"] = valid.Validator(func(str string) bool {
		ids := strings.Split(str, ",")
		if len(ids) == 0 {
			return false
		}
		return true
	})
	mt.Unlock()

	result, err := valid.ValidateStruct(commitOrder)

	if err != nil {
		return result, errors.New("参数错误")
	}

	return result, nil
}

//订单支付验证
type Payment struct {
	No        string `json:"no"`
	PayType   string `json:"pay_type" valid:"-"`
	PayMethod string `json:"pay_method" valid:"-"`
	Stage     int    `json:"stage" valid:"-"`
	Openid    string `json:"openid" valid:"-"`
}

func (payment *Payment) PaymentValidator() (bool, error) {

	if !valid.IsIn(payment.PayType, "course", "invoice", "vip") || !valid.IsIn(payment.PayMethod, "weixin_h5", "weixin_jsapi", "alipay", "hua_bai") {
		return false, errors.New("订单类型不正确/支付方式不正确")
	}

	if payment.PayMethod == "hua_bai" && !valid.IsIn(strconv.Itoa(payment.Stage), "-1", "3", "6", "12") {
		return false, errors.New("花呗分期期数不支持")
	}

	if payment.PayMethod == "weixin_jsapi" && valid.IsNull(payment.Openid) {
		return false, errors.New("openid必须")
	}

	return true, nil
}
