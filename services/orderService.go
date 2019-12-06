package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/utils"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"sync"
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
	order_type             string                 //订单类型
	user_coupon_id         int                    //优惠券id
	user_mark              string                 //用户留言
	source                 string                 //来源
	channel_uuid           int                    //渠道ID，传过来的是字符，但是自己还要取一次到id
	surface_price          float32                //总标价
	discount_price         float32                //总折扣价
	course_price           coursePrice            //课程价格信息
	available_coupon       availableCoupon        //可用的优惠券信息
	coupon_price           float32                //总优惠券价格
	db                     *BaseOrm               //数据库操作对象
	auth                   string                 //授权信息
	available_coupon_infos []models.CouponInfo    //处理后的可用优惠券信息
	orderCourse            models.OrderCourse     //处理后的课程信息
	package_id             int                    //订单套餐ID
	course_ids             []int                  //订单下的课程ID
	period_id              int                    //订单训练营期ID
	order_data             map[string]interface{} //预订单返回数据
	mt                     sync.Mutex             //针对map的读写锁
)

/**
初始化参数
*/
func initParam() {
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

	available_coupon_infos = make([]models.CouponInfo, 0)
	orderCourse = models.OrderCourse{}
	order_data = make(map[string]interface{})

	order_type = ""
	surface_price = 0
	discount_price = 0

	package_id = 0
	course_ids = make([]int, 0)
	period_id = 0

	user_coupon_id = 0
	user_mark = ""
	source = "pc"
	channel_uuid = 0
}

/**
初始化请求
*/
func initRequest(r *rest.Request, commitOrder *middlewares.CommitOrder) (int, interface{}) {
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
		db.GetDB().Table("h_edu_packages").Where("id in (?) and status = 'published'", ids).Find(&packages)
		if packages.Id == 0 {
			return 1, "未找到对应ID套餐信息"
		}

		package_id = packages.Id

		if _, err := initBaseData(packages); err != nil {
			return 1, err.Error()
		}
	} else if commitOrder.Type == "course" {
		db.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", ids).Find(&courses)
		if len(courses) == 0 {
			return 1, "未找到对应ID课程信息"
		}

		if _, err := initBaseData(courses); err != nil {
			return 1, err.Error()
		}
	} else if commitOrder.Type == "training" {
		//训练营
		db.GetDB().Table("h_edu_periods").Where("id = ? and status != 'closed'", commitOrder.PeriodId).Find(&periods)
		if periods.ID == 0 {
			return 1, "未找到对应期的信息"
		}

		db.GetDB().Table("h_edu_trainings").Where("id = ? and status = 'published'", periods.TrainingId).Find(&trainings)

		if trainings.ID == 0 {
			return 1, "未找到对应营的信息"
		}

		if _, err := initBaseData(periods); err != nil {
			return 1, err.Error()
		}
	}

	log.Info("user_coupon_id:", user_coupon_id)
	log.Info("user_coupon_id:", user_mark)
	log.Info("user_coupon_id:", source)
	log.Info("user_coupon_id:", channel_uuid)

	return 0, nil
}

/**
提交订单
*/
func (baseOrm *BaseOrm) SubmitOrder(r *rest.Request, commitOrder *middlewares.CommitOrder) (int, interface{}) {
	initParam()
	code, res := initRequest(r, commitOrder)
	if code == 1 {
		return code, res
	}

	getAvailableCoupon()
	getCourses()
	getOrderData()

	return 0, order_data
}

/**
创建订单
*/
func (baseOrm *BaseOrm) CreateOrder(r *rest.Request, commitOrder *middlewares.CommitOrder) (int, interface{}) {
	initParam()
	code, res := initRequest(r, commitOrder)
	if code == 1 {
		return code, res
	}

	if commitOrder.UserCouponId > 0 {
		user_coupon_id = commitOrder.UserCouponId
		if _, err := checkOrderCouponIsValid(); err != nil {
			return 1, err
		}
	}

	user_mark = commitOrder.UserMark
	source = commitOrder.Source
	if len(commitOrder.ChannelUuid) > 0 {
		row := db.GetDB().Table("h_market_channels").Where("uuid = ?", commitOrder.ChannelUuid).Select("id").Row()
		row.Scan(&channel_uuid)
	}

	return 0, nil
}

/**
初始化数据
*/
func initBaseData(data interface{}) (bool, error) {
	dataValue := reflect.ValueOf(data)

	_, err := checkOrderIsValid(data)
	if err != nil {
		log.Info("数据验证错误:", err.Error())
		return false, err
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			val := dataValue.Index(i).Interface().(models.Course)
			course_ids = append(course_ids, val.Id)
			//初始化价格信息
			mt.Lock()
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
			mt.Unlock()
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
			rows, err := db.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
			if err == nil {
				for rows.Next() {
					rows.Scan(&course_id)
					course_ids = append(course_ids, course_id)
				}

				db.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", course_ids).Find(&courses)
				for _, value := range courses {
					mt.Lock()
					course_price.price[value.Id] = value.Price
					course_price.discount[value.Id] = 0
					course_price.coupon[value.Id] = 0
					mt.Unlock()
				}
			}
		} else if dataValue.Type().Name() == "Period" {
			//训练营的课程，
			var course models.Course
			val := dataValue.Interface().(models.Period)
			db.GetDB().Table("h_edu_courses").Where("id = ? and status = 'published'", val.CourseId).First(&course)
			course_ids = append(course_ids, val.CourseId)
			period_id = val.ID

			//初始化价格信息
			mt.Lock()
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
			mt.Unlock()
		}
		break
	default:
		break
	}

	//计算可用优惠券的课程
	calculateAvailableCoupon()
	return true, nil
}

/**
验证订单数据的有效性
*/
func checkOrderIsValid(data interface{}) (bool, error) {
	if checkCourseIsFreeOrVipFree(data) {
		return false, utils.ReturnErrors("该订单有免费课程，无法完成下单")
	}

	if checkUserHasOrder(data) {
		return false, utils.ReturnErrors("当前课程在您未支付的订单中已经存在，无法购买")
	}

	if checkUserHasCourse(data) {
		return false, utils.ReturnErrors("课程已兑换/购买")
	}

	return true, nil
}

/**
检查订单是否包含免费课程以及会员免费学习的课程
*/
func checkCourseIsFreeOrVipFree(data interface{}) bool {
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
		rows, err := db.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
		if err == nil {
			for rows.Next() {
				rows.Scan(&course_id)
				course_ids = append(course_ids, course_id)
			}

			db.GetDB().Table("h_edu_courses").Where("id in (?) and status = 'published'", course_ids).Find(&courses)

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
		db.GetDB().Table("h_edu_courses").Where("id = ? and status = 'published'", val.CourseId).Find(&course)
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
func checkUserHasOrder(data interface{}) bool {
	switch val := data.(type) {
	case models.Course:
		var orderItems []models.OrderItemModel
		db.GetDB().Table("h_order_items").Where("course_id = ? and user_id = ? ", val.Id, user.Id).Find(&orderItems)
		if len(orderItems) > 0 {
			var order models.OrderModel
			for _, value := range orderItems {
				db.GetDB().Table("h_orders").Where("id = ? and user_id = ? and status = 0 and payment_status = 0", value.OrderId, user.Id).Find(&order)
				if order.ID > 0 {
					return true
				}
			}
		}
		return false
	case models.Package:
		var course_id = 0
		var orderItems []models.OrderItemModel
		rows, err := db.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
		if err == nil {
			for rows.Next() {
				rows.Scan(&course_id)
				course_ids = append(course_ids, course_id)
			}

			db.GetDB().Table("h_order_items").Where("course_id in (?) and user_id = ? ", course_ids, user.Id).Find(&orderItems)

			if len(orderItems) > 0 {
				var order models.OrderModel
				for _, value := range orderItems {
					db.GetDB().Table("h_orders").Where("id = ? and user_id = ? and status = 0 and payment_status = 0", value.OrderId, user.Id).Find(&order)
					if order.ID > 0 {
						return true
					}
				}
			}
		}
		return false
	case models.Period:
		var orderItems []models.OrderItemModel
		db.GetDB().Table("h_order_items").Where("course_id = ? and user_id = ?", val.CourseId, user.Id).Find(&orderItems)
		if len(orderItems) > 0 {
			var order models.OrderModel
			for _, value := range orderItems {
				db.GetDB().Table("h_orders").Where("id = ? and user_id = ? and status = 0 and payment_status = 0", value.OrderId, user.Id).Find(&order)
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
func checkUserHasCourse(data interface{}) bool {
	var user_course models.UserCourse
	switch val := data.(type) {
	case models.Course:
		db.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", val.Id, user.Id).First(&user_course)
		if user_course.Id > 0 {
			return true
		}
		return false
	case models.Package:
		rows, err := db.GetDB().Table("h_edu_package_course").Where("package_id = ?", val.Id).Select("course_id").Rows()
		if err == nil {
			var course_id = 0
			for rows.Next() {
				rows.Scan(&course_id)
				db.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", course_id, user.Id).First(&user_course)
				if user_course.Id > 0 {
					return true
				}
			}
		}
		return false
	case models.Period:
		db.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", val.CourseId, user.Id).First(&user_course)
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
func calculateAvailableCoupon() {
	if order_type == "package" {
		return
	}

	for index, value := range course_price.discount {
		if value == 0 {
			//按课程(这里展示了所有的课程)
			mt.Lock()
			available_coupon.course[index] = index
			mt.Unlock()

			//按类目
			row := db.GetDB().Table("h_edu_courses").Where("id = ? and type = 'boutique'", index).Select("category_id").Row()
			var category_id int
			err := row.Scan(&category_id)
			if err == nil {
				mt.Lock()
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
				mt.Unlock()
			}

			//按套餐(这个单独考虑，因为只有一个)
			//按课程类型(1为精品课，2为体系课)、按训练营
			row2 := db.GetDB().Table("h_edu_courses").Where("id = ?", index).Select("type").Row()
			var course_type string
			err2 := row2.Scan(&course_type)
			if err2 == nil {
				mt.Lock()
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
					db.GetDB().Table("h_edu_periods").Where("course_id = ? and status != 'closed' and sign_up_end_at > ?", index, t.Format("2006-01-02 15:04:05")).Select("id").Find(&periods)

					for _, val := range periods {
						available_coupon.training[val.ID] = index //一期只会和一个课程关联
					}
				}
				mt.Unlock()
			}
		}
	}
}

/**
返回订单预处理数据
*/
func getOrderData() {
	mt.Lock()
	order_data["surface_price"] = surface_price
	order_data["discount_price"] = discount_price
	order_data["courses"] = orderCourse
	order_data["available_coupon_infos"] = available_coupon_infos
	mt.Unlock()
}

/**
获取当前用户可用的优惠券
*/
func getAvailableCoupon() {
	//当前最低价格
	var differ = surface_price - discount_price
	var min_amount = strconv.FormatFloat(float64(differ), 'f', 2, 64)
	//至少为1块
	var max_amount = strconv.FormatFloat(float64(differ-1), 'f', 2, 64)
	var select_sql string
	var sql_str string

	couponInfos := make([]models.CouponInfo, 0)

	if order_type == "course" || order_type == "training" {
		//
		if len(available_coupon.course) == 0 {
			return
		}

		select_sql = "select uc.*, c.name as c_name, c.value as c_value, c.min_amount as c_min_amount, c.suitable as c_suitable, c.not_before as c_not_before, c.not_after as c_not_after, c.effective_day as c_effective_day from h_user_coupon as uc left join h_coupons as c on uc.coupon_id = c.id where uc.status = 0 and uc.user_id = %d and (uc.suitable = 'all'"

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

		select_sql += ") and (c.enabled = 1 and c.min_amount <=" + min_amount + " and c.value <= " + max_amount + " and c.not_before < now()) and (if (c.effective_day > 0, c.not_after > date_add(now(), interval c.effective_day day), c.not_after > now()))"
		sql_str = fmt.Sprintf(select_sql, user.Id, st, st2, st3)
	} else if order_type == "package" {
		//
		select_sql = "select uc.*, c.name as c_name, c.value as c_value, c.min_amount as c_min_amount, c.suitable as c_suitable, c.not_before as c_not_before, c.not_after as c_not_after, c.effective_day as c_effective_day from h_user_coupon as uc left join h_coupons as c on uc.coupon_id = c.id where uc.status = 0 and uc.user_id = %d and (uc.suitable = 'all' or (uc.suitable = 'package' and uc.suitable_value = " + strconv.Itoa(package_id) + ")) and (c.enabled = 1 and c.min_amount <=" + min_amount + " and c.value <= " + max_amount + " and c.not_before < now()) and (if (c.effective_day > 0, c.not_after > date_add(now(), interval c.effective_day day), c.not_after > now()))"
		sql_str = fmt.Sprintf(select_sql, user.Id)
	}

	//gorm 原生sql 读:Raw 其它操作:Exec，这里有坑，select不能用Exec
	if err := db.GetDB().Raw(sql_str).Find(&couponInfos).Error; err != nil {
		log.Info("select err is:", err.Error())
		return
	}

	//这里就不用了channel了
	if len(couponInfos) > 0 {
		var wg sync.WaitGroup
		for i := 0; i < len(couponInfos); i++ {
			wg.Add(1)
			go checkUserCouponPriceIsValid(couponInfos[i], &wg)
		}
		wg.Wait()
	}
}

/**
获取当前订单下的课程信息
*/
func getCourses() {
	db.GetDB().Table("h_edu_courses").Select("id, type, cover_picture, title, price, vip_price, discount, discount_end_at, vip_level").Where("id in (?)", course_ids).Find(&orderCourse.Courses)

	if order_type == "package" {
		var packages models.Package
		db.GetDB().Table("h_edu_packages").Where("id = " + strconv.Itoa(package_id)).First(&packages)
		orderCourse.PackageInfo = &packages
	} else if order_type == "training" {
		var period models.Period
		var training models.Training

		//join可以，但是返回值是个麻烦
		if err := db.GetDB().Table("h_edu_periods").Where("id = " + strconv.Itoa(period_id)).First(&period).Error; err != nil {
			db.GetDB().Table("h_edu_trainings").Where("id = " + strconv.Itoa(period.TrainingId)).First(&training)
			orderCourse.PeriodInfo = &period
			orderCourse.TrainingInfo = &training
		}
	}
}

/**
检查优惠券价格是否合理
*/
func checkUserCouponPriceIsValid(couponInfo models.CouponInfo, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	var all_course_price float32 = 0

	if couponInfo.CSuitable == "all" {
		for index, _ := range available_coupon.course {
			all_course_price += course_price.price[index].(float32)
		}
	} else if couponInfo.CSuitable == "category" {
		for index, value := range available_coupon.category {
			if index == couponInfo.SuitableValue {
				//这里的value为slice类型
				vals := value.([]int)
				for _, val := range vals {
					all_course_price += course_price.price[val].(float32)
				}
			}
		}
	} else if couponInfo.CSuitable == "course" {
		all_course_price += course_price.price[couponInfo.SuitableValue].(float32)
	} else if couponInfo.CSuitable == "course_type" {
		for index, value := range available_coupon.courseType {
			if index == couponInfo.SuitableValue {
				//这里的value为slice类型
				vals := value.([]int)
				for _, val := range vals {
					all_course_price += course_price.price[val].(float32)
				}
			}
		}
	} else if couponInfo.CSuitable == "training" {
		for index, value := range available_coupon.training {
			if index == couponInfo.SuitableValue {
				//这里的value为slice类型
				vals := value.([]int)
				for _, val := range vals {
					all_course_price += course_price.price[val].(float32)
				}
			}
		}
	}

	if (all_course_price-0.01) >= couponInfo.CValue && all_course_price >= couponInfo.CMinAmount {
		available_coupon_infos = append(available_coupon_infos, couponInfo)
	}
}

/**
检查当前订单的优惠券是否有效
*/
func checkOrderCouponIsValid() (bool, error) {
	//当没有优惠券的时候，也是合法的
	if user_coupon_id == 0 {
		return false, errors.New("优惠券不存在")
	}

	if order_type == "package" && discount_price > 0 {
		return false, errors.New("优惠券不合法")
	}

	if order_type == "course" && len(available_coupon.course) == 0 {
		return false, errors.New("优惠券不合法")
	}

	//
	couponInfo := models.CouponInfo{}
	if err := db.GetDB().
		Table("h_user_coupon").
		Joins("left join h_coupons on h_user_coupon.coupon_id = h_coupons.id").
		Select("h_user_coupon.*, h_coupons.name as c_name, h_coupons.value as c_value, h_coupons.min_amount as c_min_amount, h_coupons.suitable as c_suitable, h_coupons.not_before as c_not_before, h_coupons.not_after as c_not_after, h_coupons.effective_day as c_effective_day").
		Where("h_user_coupon.user_id = ? and h_user_coupon.coupon_id = ?", user.Id, user_coupon_id).
		First(&couponInfo).Error; err != nil {
		return false, errors.New("未找到用户相关优惠券信息")
	}

	if couponInfo.Status == 1 {
		return false, errors.New("优惠券已使用")
	}

	if couponInfo.Status == 2 {
		return false, errors.New("优惠券已过期")
	}

	if couponInfo.CEnabled == 0 || couponInfo.CEnabled == -1 {
		return false, errors.New("优惠券已关闭")
	}

	if couponInfo.CNotBefore.Unix() > time.Now().Unix() {
		return false, errors.New("优惠券还未生效")
	}

	if couponInfo.CEffectiveDay > 0 {
		hh := strconv.Itoa(couponInfo.CEffectiveDay*24) + "h"
		dd, _ := time.ParseDuration(hh)
		if couponInfo.CreatedAt.Add(dd).After(time.Now()) {
			return false, errors.New("优惠券已过期")
		}
	} else {
		if couponInfo.CNotAfter.Unix() < time.Now().Unix() {
			return false, errors.New("优惠券已过期")
		}
	}

	if couponInfo.CMinAmount > (surface_price-discount_price) || couponInfo.CValue > (surface_price-discount_price-1) {
		return false, errors.New("优惠券不合法")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go checkUserCouponPriceIsValid(couponInfo, &wg)
	wg.Wait()

	if len(available_coupon_infos) == 0 {
		return false, errors.New("优惠券不合法")
	}

	return true, nil
}

/**
将优惠券金额分配到当前订单下的课程里面
*/
func initOrderCouponPrice() {
	if order_type == "package" || user_coupon_id == 0 {
		return
	}

	//最终选择的优惠券，只能有一张
	var user_coupon models.CouponInfo
	if len(available_coupon_infos) > 0 {
		user_coupon = available_coupon_infos[0]
	}

	if user_coupon.ID == 0 {
		return
	}

	coupon_price = user_coupon.CValue

	coupon_course_item := make(map[int]float32)
	can_access_course_ids := make([]int, 0)
	if user_coupon.Suitable == "all" {

	}

	log.Info("coupon_course_item is:", coupon_course_item)
	log.Info("can_access_course_ids is:", can_access_course_ids)
}
