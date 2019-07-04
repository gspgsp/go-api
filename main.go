package main

import (
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"helix-edu-api/controllers/edu"
	"helix-edu-api/services"
)

func main()  {
	//初始化数据库连接实例
	new(services.BaseOrm).InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack ...)
	router, err := rest.MakeRouter(
		rest.Get("/category",new (edu.CategoryController).GetCategory),//这里传的是函数名称不需要(),只用传入方法名称
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.Println(http.ListenAndServe(":8080", nil))
}
