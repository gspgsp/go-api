package main

import (
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/controllers/edu"
	"edu_api/services"
	"edu_api/controllers/auth"
	"edu_api/middlewares"
	"regexp"
)

func main() {
	//初始化数据库连接实例
	new(services.BaseOrm).InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack ...)

	//初始化中间件
	authTokenMiddleware := new(middlewares.AuthTokenMiddleware) //或者&middlewares.AuthTokenMiddleware{}

	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {

			path := request.URL.Path

			expr := `(/login)|(/register)|(/package)|(/course[/\d+]?)|(/category)|(/chapter[/\d+]?)`
			re, _ := regexp.Compile(expr)

			all := re.FindAllString(path, -1)

			for _, item := range all {
				log.Printf("the item is:%v", string(item))
				if len(string(item)) > 0 {
					return false
				}
			}

			return true
		},
		IfTrue: authTokenMiddleware,
	})

	router, err := rest.MakeRouter(
		rest.Post("/login", new(auth.LoginController).Login),
		rest.Get("/category", new(edu.CategoryController).GetCategory),         //课程分类 这里传的是函数名称不需要(),只用传入方法名称
		rest.Get("/course/list", new(edu.CourseController).GetCourseList),      //课程列表
		rest.Get("/chapter/:id", new(edu.CourseController).GetCourseChapter),   //课程章节
		rest.Get("/package", new(edu.CourseController).GetPackageList),         //套餐列表
		rest.Get("/course/:id", new(edu.CourseController).GetCourseDetail),     //课程详情
		rest.Get("/material/:id", new(edu.MaterialController).GetMaterialList), //资料列表
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.Println(http.ListenAndServe(":8080", nil))
}
