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
		initBaseData(packages, auth, baseOrm)
	} else if commitOrder.Type == "course" {
		baseOrm.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", ids).Find(&courses)
		if len(courses) == 0 {
			return 1, "未找到对应ID课程信息"
		}
		initBaseData(courses, auth, baseOrm)
	} else if commitOrder.Type == "training" {
		//训练营
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
func initBaseData(data interface{}, auth string, baseOrm *BaseOrm) {

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
		if dataValue.Type().Name() == "Package" {
			val := dataValue.Interface().(models.Package)

			surface_price = val.Price
			formatTime, _ := utils.ParseStringTImeToStand(val.DiscountEndAt)
			if formatTime.Unix() > time.Now().Unix() {
				discount_price = val.Discount
			}
		} else if dataValue.Type().Name() == "Training" {
			//训练营的课程，
		}

		break
	default:
		break
	}

	log.Info("the surface_price is:", surface_price)
	log.Info("the discount_price is:", discount_price)
	log.Info("the type is:", dataValue.Type().Name())
}

func checkOrderIsValid(data interface{}, user models.User, baseOrm *BaseOrm) (bool, error) {
	if checkCourseIsFreeOrVipFree(data, user, baseOrm) {
		return false, utils.ReturnErrors("该订单有免费课程，无法完成下单")
	}

	if checkUserHasOrder(data, user, baseOrm) {
		return false, utils.ReturnErrors("当前课程在您未支付的订单中已经存在，无法购买")
	}

	if checkUserHasCourse(data, user, baseOrm) {
		return false, utils.ReturnErrors("课程已兑换/购买")
	}

	return true, nil
}

/**
检查订单是否包含免费课程以及会员免费学习的课程
*/
func checkCourseIsFreeOrVipFree(data interface{}, user models.User, baseOrm *BaseOrm) bool {
	switch val := data.(type) {
	case models.Course:
		if (val.Type == "free" || val.Price == 0) || (val.VipLevel == 1 && user.Level == "vip1") {
			return true
		}
		return false
	case models.Package:
		//套餐里的课程只用判断套餐是否会员免费就可以了，外是内是；外不是内不是；
		if val.VipLevel == 1 && user.Level == "vip1" {
			return true
		}

		var courses []models.Course
		var course_id = 0
		var course_ids []int
		rows, err := baseOrm.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
		if err == nil {
			for rows.Next() {
				rows.Scan(&course_id)
				course_ids = append(course_ids, course_id)
			}

			baseOrm.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", course_ids).Find(&courses)

			if len(courses) > 0 {
				for _, value := range courses {
					if value.Type == "free" {
						return true
					}
				}
			}
		}
	default:
		return false
	}

	return false
}

func checkUserHasOrder(data interface{}, user models.User, baseOrm *BaseOrm) bool {
	//switch val := data.(type) {
	//case models.Course:
	//
	//	return false
	//case models.Package:
	//
	//default:
	//	return false
	//}
	return true
}

func checkUserHasCourse(data interface{}, user models.User, baseOrm *BaseOrm) bool {

	return true
}
