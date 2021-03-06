package routes

import (
	"edu_api/controllers/auth"
	"edu_api/controllers/cart"
	"edu_api/controllers/cashier"
	"edu_api/controllers/edu"
	"edu_api/controllers/home"
	"edu_api/controllers/order"
	"edu_api/controllers/user"
	"edu_api/controllers/vip"
	"github.com/ant0ine/go-json-rest/rest"
)

/**
初始化路由
*/
func InitRoute() (rest.App, error) {

	router, err := rest.MakeRouter(
		rest.Post("/login", new(auth.LoginController).Login),
		rest.Get("/category", new(edu.CategoryController).GetCategory),                //课程分类 这里传的是函数名称不需要(),只用传入方法名称
		rest.Get("/course/list", new(edu.CourseController).GetCourseList),             //课程列表
		rest.Get("/chapter/:id", new(edu.CourseController).GetCourseChapter),          //课程章节
		rest.Get("/package", new(edu.CourseController).GetPackageList),                //套餐列表
		rest.Get("/course/:id", new(edu.CourseController).GetCourseDetail),            //课程详情
		rest.Get("/material/:id", new(edu.MaterialController).GetMaterialList),        //资料列表
		rest.Get("/lecture/:id", new(user.UserController).GetLecturerList),            //讲师列表
		rest.Get("/review/:id", new(edu.CourseController).GetCourseReview),            //评价列表
		rest.Get("/recommend/:id", new(edu.CourseController).GetRecommendCourse),      //推荐课程
		rest.Get("/play/:id/:lesion_id", new(edu.PlayController).GetPlayList),         //视频播放
		rest.Post("/learn", new(edu.PlayController).PutCourseLearn),                   //视频观看记录
		rest.Post("/remark/:id", new(edu.RemarkController).StoreRemark),               //创建评价
		rest.Get("/try/:id", new(edu.CourseController).GetTrySeeList),                 //获取试看列表
		rest.Get("/compose/:id", new(edu.PackageController).GetComposePackage),        //获取组合套餐
		rest.Get("/package/:id", new(edu.PackageController).GetPackageDetail),         //套餐详情
		rest.Get("/exam/rolls/:id", new(edu.ExamController).GetExamRollTopicList),     //获取题库作业列表
		rest.Get("/exam/roll/:id", new(edu.ExamController).GetExamRollTopicInfo),      //获取题库作业详情
		rest.Post("/exam/roll/:id", new(edu.ExamController).StoreTopicAnswer),         //提交答案
		rest.Get("/exam/roll/:id/:course_id", new(edu.ExamController).GetTopicAnswer), //答案解析
		rest.Get("/vip/:id", new(vip.VipController).GetVipInfo),                       //VIP信息
		rest.Post("/vip/add", new(vip.VipController).CreateVipOrder),                  //创建会员订单
		rest.Delete("/vip/:id", new(vip.VipController).DeleteVipOrder),                //取消会员订单
		rest.Get("/notice", new(home.NoticeController).GetNotice),                     //公告信息
		rest.Get("/slide", new(home.SlideController).GetSlide),                        //轮播信息
		rest.Post("/cart", new(cart.CartController).AddCartInfo),                      //添加购物车
		rest.Get("/cart", new(cart.CartController).GetCartList),                       //购物车列表
		rest.Delete("/cart/:id", new(cart.CartController).DelCart),                    //删除购物车
		rest.Post("/order/submit", new(order.OrderController).SubmitOrder),            //提交订单
		rest.Post("/order/create", new(order.OrderController).CreateOrder),            //创建订单
		rest.Post("/cashier/payment", new(cashier.CashierController).Payment),         //生成支付信息
		rest.Post("/pay_notify/:type", new(cashier.CashierController).PayNotify),      //支付信息异步通知统一接口
	)

	if err != nil {
		return nil, err
	}

	return router, nil

}
