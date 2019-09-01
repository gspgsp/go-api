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
	"edu_api/hook"
	"time"
	"github.com/lestrrat-go/file-rotatelogs"
)

//初始化日志操作,全局有效
func initLog() {
	//测试logrus日志包，这个包有个依赖，golang.org/x/sys
	//log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	//设置最低loglevel
	//log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	//file, _ := os.OpenFile("./src/edu_api/log/request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//log.SetOutput(file)
	configLocalFilesystemLogger("./src/edu_api/log/request")
}

//切割日志和清理过期日志
func configLocalFilesystemLogger(filePath string) {
	writer, err := rotatelogs.New(
		filePath+".%Y%m%d.log",                    //%Y%m%d%H%M"，日志分割时间
		rotatelogs.WithLinkName(filePath),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*7*24),     // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)
	if err != nil {
		log.Fatal("Init log failed, err:", err)
	}
	log.SetOutput(writer)
}

func main() {

	//初始化日志操作
	initLog()

	//所有的日志里面通过自定义hook函数追加一个，字符串
	log.AddHook(hook.NewTraceInfoHook("最终解释权归GJH"))

	//启动拂去提示
	log.Info("启动服务...")

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

	//监听端口提示
	log.Info("监听端口:", utils.SERVER_PORT)
	log.Error(http.ListenAndServe(":"+utils.SERVER_PORT, nil))
}
