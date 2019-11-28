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
	var periods models.Period
	var trainings models.Training
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
		baseOrm.GetDB().Table("h_edu_periods").Where("id = ? and status != 'closed'", commitOrder.PeriodId).Find(&periods)
		if periods.ID == 0 {
			return 1, "未找到对应期的信息"
		}

		baseOrm.GetDB().Table("h_edu_trainings").Where("id = ? and status = 'published'", periods.TrainingId).Find(&trainings)

		if trainings.ID == 0 {
			return 1, "未找到对应营的信息"
		}
		initBaseData(periods, auth, baseOrm)
	}

	return 0, "ok"
}

/**
初始化数据
*/
func initBaseData(data interface{}, auth string, baseOrm *BaseOrm) {

	res, err := checkOrderIsValid(data, user, baseOrm)
	if res && err != nil {
		log.Info("数据验证错误:", err.Error())
	}

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
		} else if dataValue.Type().Name() == "Period" {
			//训练营的课程，
			//val := dataValue.Interface().(models.Period)

		}
		break
	default:
		break
	}

	log.Info("the surface_price is:", surface_price)
	log.Info("the discount_price is:", discount_price)
	log.Info("the type is:", dataValue.Type().Name())
}

/**
验证订单数据的有效性
*/
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
	case models.Period:
		var course models.Course
		baseOrm.GetDB().Table("h_edu_courses").Where("id = ? and status = 'published'", val.CourseId).Find(&course)
		if course.Type == "free" || course.Price == 0 || (course.VipLevel == 1 && user.Level == "vip1") {
			return true
		}
	default:
		return false
	}

	return false
}

/**
检查是否有未处理的订单，其实这个表结构设计有问题，就是没有设置order_items的状态，每次都要到父表查询状态
*/
func checkUserHasOrder(data interface{}, user models.User, baseOrm *BaseOrm) bool {
	switch val := data.(type) {
	case models.Course:
		var orderItems []models.OrderItemModel
		baseOrm.GetDB().Table("h_order_items").Where("course_id = ? and user_id = ? ", val.Id, user.Id).Find(&orderItems)
		if len(orderItems) > 0 {
			var order models.OrderModel
			for _, value := range orderItems {
				baseOrm.GetDB().Table("h_orders").Where("id = ? and user_id = ? and status = 0 and payment_status = 0", value.OrderId, user.Id).Find(&order)
				if order.ID > 0 {
					return true
				}
			}
		}
		return false
	case models.Package:
		var course_id = 0
		var course_ids []int
		var orderItems []models.OrderItemModel
		rows, err := baseOrm.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
		if err == nil {
			for rows.Next() {
				rows.Scan(&course_id)
				course_ids = append(course_ids, course_id)
			}

			baseOrm.GetDB().Table("h_order_items").Where("course_id in (?) and user_id = ? ", course_ids, user.Id).Find(&orderItems)

			if len(orderItems) > 0 {
				var order models.OrderModel
				for _, value := range orderItems {
					baseOrm.GetDB().Table("h_orders").Where("id = ? and user_id = ? and status = 0 and payment_status = 0", value.OrderId, user.Id).Find(&order)
					if order.ID > 0 {
						return true
					}
				}
			}
		}
		return false
	case models.Period:
		var orderItems []models.OrderItemModel
		baseOrm.GetDB().Table("h_order_items").Where("course_id = ? and user_id = ?", val.CourseId, user.Id).Find(&orderItems)
		if len(orderItems) > 0 {
			var order models.OrderModel
			for _, value := range orderItems {
				baseOrm.GetDB().Table("h_orders").Where("id = ? and user_id = ? and status = 0 and payment_status = 0", value.OrderId, user.Id).Find(&order)
				if order.ID > 0 {
					return true
				}
			}
		}
		return false
	default:
		return false
	}
}

/**
检查用户是否已经有指定的课程
*/
func checkUserHasCourse(data interface{}, user models.User, baseOrm *BaseOrm) bool {
	var user_course models.UserCourse
	switch val := data.(type) {
	case models.Course:
		baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", val.Id, user.Id).First(&user_course)
		if user_course.Id > 0 {
			return true
		}
		return false
	case models.Package:
		rows, err := baseOrm.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
		if err == nil {
			var course_id = 0
			for rows.Next() {
				rows.Scan(&course_id)
				baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", course_id, user.Id).First(&user_course)
				if user_course.Id > 0 {
					return true
				}
			}
		}
		return false
	case models.Period:
		baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", val.CourseId, user.Id).First(&user_course)
		if user_course.Id > 0 {
			return true
		}
		return false
	default:
		return false
	}
}
