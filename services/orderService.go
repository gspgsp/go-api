package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/utils"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

var (
	surface_price  float32
	discount_price float32
)

/**
提交订单
*/
func (baseOrm *BaseOrm) SubmitOrder(r *rest.Request, commitOrder *middlewares.CommitOrder) (int, interface{}) {

	var courses []models.Course
	var packages models.Package
	var ids []string
	auth := r.Header.Get("Authorization")
	ids = strings.Split(commitOrder.Ids, ",")

	if commitOrder.Type == "package" {
		baseOrm.GetDB().Table("h_edu_packages").Where("id in (?) and status = 'published'", ids).Find(&packages)
		if packages.Id == 0 {
			return 1, "未找到对应ID套餐信息"
		}
		initBaseData(packages, auth)
	} else {
		baseOrm.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", ids).Find(&courses)
		if len(courses) == 0 {
			return 1, "未找到对应ID课程信息"
		}
		initBaseData(courses, auth)
	}

	//'amount' => $this->formatFloat($this->amount),
	//'discount_amount' => $this->formatFloat($this->discount_amount),
	//'goods' => $this->order_type == 'package' ? $this->package : $this->courses->values(),
	//	'coupons' => $this->getAvailableCoupons(),

	return 0, "ok"
}

/**
初始化数据
*/
func initBaseData(data interface{}, auth string) float32 {

	dataValue := reflect.ValueOf(data)
	user = GetUserInfo(auth)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			val := dataValue.Index(i).Interface().(models.Course)
			if user.Level != "vip1" {
				surface_price += val.Price

				formatTime, _ := utils.ParseStringTImeToStand(val.DiscountEndAt)
				if formatTime.Unix() > time.Now().Unix() {
					discount_price += val.Discount
				}
			} else {
				if val.VipLevel == 0 && val.VipPrice > 0 {
					surface_price += val.VipPrice
				} else {
					surface_price += val.Price
				}
			}
		}
		break
	case reflect.Struct:
		val := dataValue.Interface().(models.Package)

		surface_price = val.Price
		formatTime, _ := utils.ParseStringTImeToStand(val.DiscountEndAt)
		if formatTime.Unix() > time.Now().Unix() {
			discount_price = val.Discount
		}

		break
	default:
		break
	}

	log.Info("the surface_price is:", surface_price)
	log.Info("the discount_price is:", discount_price)

	return 0.23
}
