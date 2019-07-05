package main

import (
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/controllers/edu"
	"edu_api/services"
)

func main()  {
	//初始化数据库连接实例
	new(services.BaseOrm).InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack ...)
	router, err := rest.MakeRouter(
		rest.Get("/category", new (edu.CategoryController).GetCategory),//这里传的是函数名称不需要(),只用传入方法名称
		rest.Get("/course", new(edu.CourseController).GetCourseList),
		rest.Get("/package", new(edu.CourseController).GetPackageList),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.Println(http.ListenAndServe(":8080", nil))
}
