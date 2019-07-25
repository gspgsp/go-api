package main

import (
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/controllers/edu"
	"edu_api/services"
	"edu_api/controllers/auth"
)

func main()  {
	//初始化数据库连接实例
	new(services.BaseOrm).InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack ...)

	//tokenAuthMiddleware := rest.Middleware()
	//
	/*api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			var arr = []string{
				"/login","/register",
			}
			for _, item := range arr {
				if item == request.URL.Path {
					return false
				}
			}
			return true
		},
		IfTrue: tokenAuthMiddleware,
	})*/

	router, err := rest.MakeRouter(
		rest.Post("/login", new(auth.LoginController).Login),
		rest.Get("/category", new(edu.CategoryController).GetCategory),//课程分类 这里传的是函数名称不需要(),只用传入方法名称
		rest.Get("/course", new(edu.CourseController).GetCourseList),//课程列表
		rest.Get("/package", new(edu.CourseController).GetPackageList),//套餐列表
		rest.Get("/course/:id", new(edu.CourseController).GetCourseDetail),//课程详情
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.Println(http.ListenAndServe(":8080", nil))
}
