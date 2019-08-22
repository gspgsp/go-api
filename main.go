package main

import (
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/services"
	"edu_api/middlewares"
	"edu_api/routes"
	"edu_api/utils"
	"github.com/garyburd/redigo/redis"
)

func main() {
	//初始化数据库连接实例
	new(services.BaseOrm).InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack ...)

	//路由中间件
	api.Use(middlewares.IfMiddleware())

	//路由信息
	router, err := routes.InitRoute()

	if err != nil {
		log.Fatal(err)
	}

	redis.DialConnectTimeout(1)

	api.SetApp(router)

	http.Handle(utils.API_PREFIX+"/", http.StripPrefix(utils.API_PREFIX, api.MakeHandler()))

	log.Println(http.ListenAndServe(":"+utils.SERVER_PORT, nil))
}
