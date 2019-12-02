package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/utils"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

/**
课程价格信息
*/
type coursePrice struct {
	price    map[int]interface{}
	discount map[int]interface{}
	coupon   map[int]interface{}
	payment  map[int]interface{}
}

/**
课程优惠券
*/
type availableCoupon struct {
	course     map[int]int
	category   map[int]interface{}
	courseType map[int]interface{}
	training   map[int]interface{}
}

var (
	order_type       string          //订单类型
	surface_price    float32         //总标价
	discount_price   float32         //总折扣价
	course_price     coursePrice     //课程价格信息
	available_coupon availableCoupon //可用的优惠券信息
	db               *BaseOrm        //数据库操作对象
	auth             string          //授权信息
)

/**
map初始化，用于存储可用的优惠券信息，按照优惠券分类来区分
*/
func Init() {
	course_price.price = make(map[int]interface{})
	course_price.discount = make(map[int]interface{})
	course_price.coupon = make(map[int]interface{})
	course_price.payment = make(map[int]interface{})

	available_coupon.course = make(map[int]int)
	available_coupon.category = make(map[int]interface{})
	available_coupon.courseType = make(map[int]interface{})
	available_coupon.training = make(map[int]interface{})

	//数据库初始化
	db = new(BaseOrm)
	db.InitDB()
}

/**
提交订单
*/
func (baseOrm *BaseOrm) SubmitOrder(r *rest.Request, commitOrder *middlewares.CommitOrder) (int, interface{}) {

	Init()

	log.Info("the db is:", db.GetDB())
	var courses []models.Course
	var packages models.Package
	var periods models.Period
	var trainings models.Training
	var ids []string
	auth = r.Header.Get("Authorization")
	user = GetUserInfo(auth)
	ids = strings.Split(commitOrder.Ids, ",")
	order_type = commitOrder.Type

	if commitOrder.Type == "package" {
		baseOrm.GetDB().Table("h_edu_packages").Where("id in (?) and status = 'published'", ids).Find(&packages)
		if packages.Id == 0 {
			return 1, "未找到对应ID套餐信息"
		}

		if _, err := initBaseData(packages, baseOrm); err != nil {
			return 1, err.Error()
		}
	} else if commitOrder.Type == "course" {
		baseOrm.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", ids).Find(&courses)
		if len(courses) == 0 {
			return 1, "未找到对应ID课程信息"
		}

		if _, err := initBaseData(courses, baseOrm); err != nil {
			return 1, err.Error()
		}
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

		if _, err := initBaseData(periods, baseOrm); err != nil {
			return 1, err.Error()
		}
	}

	//获取返回数据
	getData()

	return 0, "ok"
}

/**
初始化数据
*/
func initBaseData(data interface{}, baseOrm *BaseOrm) (bool, error) {
	dataValue := reflect.ValueOf(data)

	_, err := checkOrderIsValid(data, baseOrm)
	if err != nil {
		log.Info("数据验证错误:", err.Error())
		return false, err
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			val := dataValue.Index(i).Interface().(models.Course)

			//初始化价格信息
			course_price.price[val.Id] = val.Price
			course_price.discount[val.Id] = 0
			course_price.coupon[val.Id] = 0

			if user.Level != "vip1" {
				surface_price += val.Price

				formatTime, _ := utils.ParseStringTImeToStand(val.DiscountEndAt)
				if formatTime.Unix() > time.Now().Unix() {
					discount_price += val.Discount
					course_price.discount[val.Id] = val.Discount
				}
			} else {
				if val.VipLevel == 0 && val.VipPrice > 0 { //会员非免费，同时又会员价
					surface_price += val.VipPrice
					course_price.price[val.Id] = val.VipPrice
				} else { //会员非免费，木有会员价
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
				for _, value := range courses {
					course_price.price[value.Id] = value.Price
					course_price.discount[value.Id] = 0
					course_price.coupon[value.Id] = 0
				}
			}
		} else if dataValue.Type().Name() == "Period" {
			//训练营的课程，
			var course models.Course
			val := dataValue.Interface().(models.Period)
			baseOrm.GetDB().Table("h_edu_courses").Where("id = ? and status = 'published'", val.CourseId).First(&course)

			//初始化价格信息
			course_price.price[course.Id] = course.Price
			course_price.discount[course.Id] = 0
			course_price.coupon[course.Id] = 0

			if user.Level != "vip1" {
				surface_price = course.Price

				formatTime, _ := utils.ParseStringTImeToStand(course.DiscountEndAt)
				if formatTime.Unix() > time.Now().Unix() {
					discount_price = course.Discount
					course_price.discount[course.Id] = course.Discount
				}
			} else {
				if course.VipLevel == 0 && course.VipPrice > 0 {
					surface_price = course.VipPrice
					course_price.price[course.Id] = course.VipPrice
				} else {
					surface_price = course.Price
				}
			}
		}
		break
	default:
		break
	}

	//计算可用优惠券的课程
	calculateAvailableCoupon(baseOrm)

	log.Info("course is:", course_price)
	log.Info("coupon is:", available_coupon.training)
	return true, nil
}

/**
验证订单数据的有效性
*/
func checkOrderIsValid(data interface{}, baseOrm *BaseOrm) (bool, error) {
	if checkCourseIsFreeOrVipFree(data, baseOrm) {
		return false, utils.ReturnErrors("该订单有免费课程，无法完成下单")
	}

	if checkUserHasOrder(data, baseOrm) {
		return false, utils.ReturnErrors("当前课程在您未支付的订单中已经存在，无法购买")
	}

	if checkUserHasCourse(data, baseOrm) {
		return false, utils.ReturnErrors("课程已兑换/购买")
	}

	return true, nil
}

/**
检查订单是否包含免费课程以及会员免费学习的课程
*/
func checkCourseIsFreeOrVipFree(data interface{}, baseOrm *BaseOrm) bool {
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
func checkUserHasOrder(data interface{}, baseOrm *BaseOrm) bool {
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
func checkUserHasCourse(data interface{}, baseOrm *BaseOrm) bool {
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

/**
计算价格为0的课程，计算可用优惠券的课程
*/
func calculateAvailableCoupon(baseOrm *BaseOrm) {
	if order_type == "package" {
		return
	}

	for index, value := range course_price.discount {
		if value == 0 {
			//按课程(这里展示了所有的课程)
			available_coupon.course[index] = index

			//按类目
			row := baseOrm.GetDB().Table("h_edu_courses").Where("id = ? and type = 'boutique'", index).Select("category_id").Row()
			var category_id int
			err := row.Scan(&category_id)
			if err == nil {
				var c_ids []int
				switch val := available_coupon.category[category_id].(type) {
				case []int:
					val = append(val, index)
					c_ids = append(c_ids, val...)
					available_coupon.category[category_id] = c_ids
					break
				default:
					c_ids = append(c_ids, index)
					available_coupon.category[category_id] = c_ids
				}
			}

			//按套餐(这个单独考虑，因为只有一个)
			//按课程类型(1为精品课，2为体系课)、按训练营
			row2 := baseOrm.GetDB().Table("h_edu_courses").Where("id = ?", index).Select("type").Row()
			var course_type string
			err2 := row2.Scan(&course_type)
			if err2 == nil {
				if course_type == "boutique" {
					var c_ids []int
					switch val := available_coupon.courseType[1].(type) {
					case []int:
						val = append(val, index)
						c_ids = append(c_ids, val...)
						available_coupon.courseType[1] = c_ids
						break
					default:
						c_ids = append(c_ids, index)
						available_coupon.courseType[1] = c_ids
					}
				} else if course_type == "class" {
					var c_ids []int
					switch val := available_coupon.courseType[2].(type) {
					case []int:
						val = append(val, index)
						c_ids = append(c_ids, val...)
						available_coupon.courseType[2] = c_ids
						break
					default:
						c_ids = append(c_ids, index)
						available_coupon.courseType[2] = c_ids
					}
				} else if course_type == "training" {
					var periods []models.Period
					t := time.Now()
					baseOrm.GetDB().Table("h_edu_periods").Where("course_id = ? and status != 'closed' and sign_up_end_at > ?", index, t.Format("2006-01-02 15:04:05")).Select("id").Find(&periods)

					for _, val := range periods {
						available_coupon.training[val.ID] = index //一期只会和一个课程关联
					}
				}
			}
		}
	}
}

/**
返回数据
*/
func getData() {
	getAvailableCoupon()
}

/**
获取商品表面价格
*/
func getSurfacePrice() float32 {
	return surface_price
}

/**
获取总的折扣价格
*/
func getDiscountPrice() float32 {
	return discount_price
}

/**
获取当前用户可用的优惠券
*/
func getAvailableCoupon() {
	if order_type == "course" {
		//
		if len(available_coupon.course) == 0 {
			return
		}

		//select_sql := "select uc.id as user_coupon_id, uc.status as user_coupon_status, c.* from h_user_coupon as uc left join h_coupons as c on uc.coupon_id = c.id where uc.status = 0 and uc.user_id = %d and (uc.suitable = 'all' or (uc.suitable = 'category' and uc.suitable_value in (%s)) or (uc.suitable = 'course' and uc.suitable_value in (%s)) or (uc.suitable = 'course_type' and uc.suitable_value in (%s)))"
		select_sql := "select uc.id as user_coupon_id, uc.status as user_coupon_status, c.* from h_user_coupon as uc left join h_coupons as c on uc.coupon_id = c.id where uc.status = 0 and uc.user_id = %d and (uc.suitable = 'all'"

		keys := make([]int, 0, len(available_coupon.category))
		for k := range available_coupon.category {
			keys = append(keys, k)
		}

		keys2 := make([]int, 0, len(available_coupon.course))
		for k := range available_coupon.course {
			keys2 = append(keys2, k)
		}

		keys3 := make([]int, 0, len(available_coupon.courseType))
		for k := range available_coupon.courseType {
			keys3 = append(keys3, k)
		}

		st := strings.Replace(strings.Trim(fmt.Sprint(keys), "[]"), " ", ",", -1)
		st2 := strings.Replace(strings.Trim(fmt.Sprint(keys2), "[]"), " ", ",", -1)
		st3 := strings.Replace(strings.Trim(fmt.Sprint(keys3), "[]"), " ", ",", -1)

		if len(st) > 0 {
			select_sql += "or (uc.suitable = 'category' and uc.suitable_value in (%s))"
		}

		if len(st2) > 0 {
			select_sql += "or (uc.suitable = 'course' and uc.suitable_value in (%s))"
		}

		if len(st3) > 0 {
			select_sql += "or (uc.suitable = 'course_type' and uc.suitable_value in (%s))"
		}

		select_sql += ")"

		sql_str := fmt.Sprintf(select_sql, user.Id, st, st2, st3)
		rows, _ := db.GetDB().Exec(sql_str).Rows()

		log.Info("the rows is:", rows)
	} else if order_type == "package" {
		//

	}
}

/**
获取当前订单下的课程信息
*/
func getCourses() {

}
