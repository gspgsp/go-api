package routes

import (
	"edu_api/controllers/auth"
	"edu_api/controllers/edu"
	"edu_api/controllers/user"
	"github.com/ant0ine/go-json-rest/rest"
)

/**
初始化路由
*/
func InitRoute() (rest.App, error) {

	router, err := rest.MakeRouter(
		rest.Post("/login", new(auth.LoginController).Login),
		rest.Get("/category", new(edu.CategoryController).GetCategory),           //课程分类 这里传的是函数名称不需要(),只用传入方法名称
		rest.Get("/course/list", new(edu.CourseController).GetCourseList),        //课程列表
		rest.Get("/chapter/:id", new(edu.CourseController).GetCourseChapter),     //课程章节
		rest.Get("/package", new(edu.CourseController).GetPackageList),           //套餐列表
		rest.Get("/course/:id", new(edu.CourseController).GetCourseDetail),       //课程详情
		rest.Get("/material/:id", new(edu.MaterialController).GetMaterialList),   //资料列表
		rest.Get("/lecture/:id", new(user.UserController).GetLecturerList),       //讲师列表
		rest.Get("/review/:id", new(edu.CourseController).GetCourseReview),       //评价列表
		rest.Get("/recommend/:id", new(edu.CourseController).GetRecommendCourse), //推荐课程
		rest.Get("/play/:id/:lesion_id", new(edu.PlayController).GetPlayList),    //视频播放
		rest.Post("/learn", new(edu.PlayController).PutCourseLearn),              //视频观看记录
		rest.Post("/remark/:id", new(edu.RemarkController).StoreRemark),          //创建评价
		rest.Get("/try/:id", new(edu.CourseController).GetTrySeeList),            //获取试看列表
		rest.Get("/compose/:id", new(edu.PackageController).GetComposePackage),   //获取组合套餐
	)

	if err != nil {
		return nil, err
	}

	return router, nil

}
