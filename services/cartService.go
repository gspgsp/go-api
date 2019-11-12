package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"strconv"
)

/**
添加购物车
*/
func (baseOrm *BaseOrm) AddCartInfo(r *rest.Request, addCart *middlewares.AddCart) (int, interface{}) {
	course_id := addCart.CourseId

	//验证
	var course models.Course
	if err := baseOrm.GetDB().Table("h_edu_courses").Where("id = ? and status = 'published'", course_id).First(&course).Error; err != nil {
		log.Info("获取数据错误:" + err.Error())
		return 1, err.Error()
	}

	if course.Type == "free" || course.Price == 0 {
		log.Info("免费课程无法加入购物车")
		return 1, "免费课程无法加入购物车"
	}

	user = GetUserInfo(r.Header.Get("Authorization"))
	var orderItems []models.OrderItemModel
	baseOrm.GetDB().Table("h_order_items").Where("course_id = ? and user_id = ?", course_id, user.Id).Select("order_id").Find(&orderItems)

	if len(orderItems) > 0 {
		var ids []int64
		var orders []models.OrderModel
		for _, value := range orderItems {
			ids = append(ids, value.OrderId)
		}

		baseOrm.GetDB().Table("h_orders").Where("id in (?) and (status = 0 or status = 1) and (payment_status = 0 or payment_status = 1)", ids).Find(&orders)
		if len(orders) > 0 {
			log.Info("当前课程在您未/已支付的订单中已经存在，无法再次购买")
			return 1, "当前课程在您未/已支付的订单中已经存在，无法再次购买"
		}
	}

	var userCourse models.UserCourse
	baseOrm.GetDB().Table("h_user_course").Where("user_id = ? and course_id = ? ", user.Id, course_id).First(&userCourse)

	if userCourse.Id > 0 {
		log.Info("您当前订单中已经有课程在学习计划中，无法再次购买")
		return 1, "您当前订单中已经有课程在学习计划中，无法再次购买"
	}

	var cart models.Cart
	baseOrm.GetDB().Table("h_carts").Where("user_id = ? and course_id = ?", user.Id, course_id).First(&cart)

	if cart.ID > 0 {
		log.Info("商品已经在购物车内")
		return 1, "商品已经在购物车内"
	}

	//添加
	if err := baseOrm.GetDB().Table("h_carts").Exec("insert into h_carts (`user_id`, `course_id`) values(?,?)", user.Id, course_id).Error; err != nil {
		log.Info("添加失败")
		return 1, "添加失败"
	}

	return 0, "添加成功"
}

/**
购物车列表
*/
func (baseOrm *BaseOrm) GetCartList(r *rest.Request) (int, interface{}) {
	user = GetUserInfo(r.Header.Get("Authorization"))

	var cartList []models.Cart

	if err := baseOrm.GetDB().Table("h_carts").Where("user_id = ?", user.Id).Find(&cartList).Error; err != nil {
		log.Info("获取数据错误:" + err.Error())
		return 1, err.Error()
	}

	for index, value := range cartList {
		baseOrm.GetDB().Table("h_edu_courses").Where("id = ?", value.CourseId).First(&cartList[index].Course)
	}

	return 0, cartList
}

/**
删除购物车
*/
func (baseOrm *BaseOrm) DelCart(r *rest.Request) (int, interface{}) {
	user = GetUserInfo(r.Header.Get("Authorization"))
	id, _ := strconv.Atoi(r.PathParam("id"))
	cart := models.Cart{ID: id}

	if err := baseOrm.GetDB().Table("h_carts").Where("user_id = ?", user.Id).Delete(&cart).Error; err != nil {
		log.Info("删除数据错误:" + err.Error())
		return 1, err.Error()
	}

	return 0, "删除成功"
}
