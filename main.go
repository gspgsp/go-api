package main

import (
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/services"
	"edu_api/middlewares"
	"edu_api/routes"
	"edu_api/utils"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"os"
	"edu_api/hook"
)

//初始化日志操作,全局有效
func initLog() {
	//测试logrus日志包，这个包有个依赖，golang.org/x/sys
	//log.SetFormatter(&log.TextFormatter{DisableTimestamp: true}),这个条件有bug，加上以后会导致下面的json格式输出有问题

	//设置最低loglevel
	//log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.JSONFormatter{})
	file, _ := os.OpenFile("./src/edu_api/log/request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)
}

func main() {

	//初始化日志操作
	initLog()

	//所有的日志里面通过自定义hook函数追加一个，字符串
	log.AddHook(hook.NewTraceIdHook("1qaz2wsx"))

	log.Info("我是测试")
	log.WithFields(log.Fields{
		"age":  12,
		"name": "xiaoming",
		"sex":  1,
	}).Info("小明来了")

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

	log.Printf("we are now at:%s", utils.SERVER_PORT)
	log.Error(http.ListenAndServe(":"+utils.SERVER_PORT, nil))
}
