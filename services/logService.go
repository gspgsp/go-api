package services

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"path"
	"time"
)

type Log struct {
	LogPath string
	LogName string
}

func (initLog *Log) InitLog() {
	//测试logrus日志包，这个包有个依赖，golang.org/x/sys
	//log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	//设置最低loglevel
	//log.SetLevel(log.InfoLevel)
	//file, _ := os.OpenFile("./src/edu_api/log/request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//log.SetOutput(file)

	var logger = log.New()

	//格式换时间输出格式
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	//日志保存到本地-日志路劲
	file := path.Join(initLog.LogPath, initLog.LogName)
	configLocalFilesystemLogger(file)

	//所有的日志里面通过自定义hook函数追加一个，字符串，定义es的hook将日志存储到es
	//log.AddHook(hook.NewTraceInfoHook("最终解释权归GJH"))
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://47.97.126.65:9200"))
	if err != nil {
		log.Info("init es err:", err.Error())
	}

	hook, err := elogrus.NewElasticHook(client, "127.0.0.1", log.DebugLevel, "mylog")
	if err != nil {
		log.Info("init hook err:", err.Error())
	}
	logger.Hooks.Add(hook)
	logger.WithFields(log.Fields{
		"name": "gjh",
		"age":  42,
		"host": "127.0.0.1",
	})
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
