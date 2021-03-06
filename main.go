package main

import (
	"edu_api/middlewares"
	"edu_api/routes"
	"edu_api/services"
	"edu_api/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/iGoogle-ink/gopay"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.Info("GoPAy version is:", gopay.Version)
	log.Info("值的纪念的日子，安全抵达上海半个多月:", "2020年4月2日9:36分")
	//初始化日志操作
	(&(services.Log{utils.LOG_PATH, utils.LOG_NAME})).InitLog()

	//初始化数据库连接实例
	new(services.BaseOrm).InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	//路由中间件
	api.Use(middlewares.IfMiddleware())

	//路由信息
	router, err := routes.InitRoute()

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle(utils.API_PREFIX+"/", http.StripPrefix(utils.API_PREFIX, api.MakeHandler()))

	//监听端口提示
	log.Info("监听端口:", utils.SERVER_PORT)
	log.Error(http.ListenAndServe(":"+utils.SERVER_PORT, nil))
}
