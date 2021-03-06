package services

import (
	"bytes"
	"edu_api/models"
	"edu_api/utils"
	"encoding/json"
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

/**
基本数据库服务类
*/
type BaseOrm struct {
	DB *gorm.DB
}

var baseDb *gorm.DB

/**
初始化数据库
*/
func (that *BaseOrm) InitDB() {
	var err error

	//获取数据库配置对象
	jsonStruct := utils.InitJsonStruct()

	v := utils.ConfigType{}

	err = jsonStruct.Load("D:/gopath/src/edu_api/config/database.json", &v)

	if err != nil {
		log.Println("parse db config error!", err)
		return
	}

	that.DB, err = gorm.Open(v.Mysql.DbConnect, v.Mysql.DbUsername+":"+v.Mysql.DbPassword+"@tcp("+v.Mysql.DbHost+":"+v.Mysql.DbPort+")/"+v.Mysql.DbDatabase+"?charset=utf8&parseTime=true&loc=Local")

	if err != nil {
		log.Println("init db error!", err)
		return
	}

	that.DB.SingularTable(true)
	that.DB.DB().SetMaxIdleConns(1000)
	that.DB.DB().SetMaxOpenConns(2600) //2核8G Xeon 最大连接数26850
	that.DB.DB().SetConnMaxLifetime(300 * time.Second)
	that.DB.LogMode(true) //决定在终端是否打印sql，开发环境下设置为true

	baseDb = that.DB
}

/**
获取数据库信息
*/
func (that *BaseOrm) GetDB() (DB *gorm.DB) {

	if baseDb != nil {
		DB = baseDb
	} else {
		log.Println("init db error")
		return
	}

	return
}

/**
预处理数据
*/
func Trees(list interface{}) interface{} {
	data := formatDatas(list)
	result := formatCores(0, data)

	return result
}

/**
格式化数据
*/
func formatDatas(list interface{}) interface{} {

	switch value := list.(type) {
	case []models.Category:

		data := make(map[int]map[int]models.Category)

		for _, v := range value {
			id := int(v.Id)
			parentId := int(v.ParentId)

			if _, ok := data[parentId]; !ok {
				//如果parent_id相同，那么就申请一个空的map(初始化后才能使用)
				data[parentId] = make(map[int]models.Category)
			}

			//将相同parent_id的数据放到统一map下
			data[parentId][id] = v
		}

		return data
	case []models.Chapter:

		data := make(map[int]map[int]models.Chapter)

		for _, v := range value {
			id := int(v.Id)
			parentId := int(v.ParentId)

			if _, ok := data[parentId]; !ok {
				data[parentId] = make(map[int]models.Chapter)
			}

			data[parentId][id] = v
		}

		return data
	}

	return nil
}

/**
格式化数据核心方法
*/
func formatCores(index int, data interface{}) interface{} {

	switch value := data.(type) {
	case map[int]map[int]models.Category:
		//初始化一个大小为0的slice
		tmp := make([]models.Category, 0)

		for id, item := range value[index] {
			if value[id] != nil {
				//用当前id去取值，parent_id会和id对应起来，注意interface{}==>指定类型的转换
				item.Children = formatCores(id, value).([]models.Category)
			}

			tmp = append(tmp, item)
		}

		return tmp
	case map[int]map[int]models.Chapter:
		tmp := make([]models.Chapter, 0)

		for id, item := range value[index] {
			if value[id] != nil {
				item.Children = formatCores(id, value).([]models.Chapter)
			}

			tmp = append(tmp, item)
		}

		return tmp
	}

	return nil
}

/**
格式化时间
*/
func FormatTime(time int64) (format_time int64, err error) {
	if time > 0 {
		format_time, err := valid.ToInt(fmt.Sprintf("%v", time/60))
		if err != nil {
			return 0, err
		}
		return format_time, nil
	}
	return 0, nil
}

func FormatTimeToChinese(time int64) (format_time string) {
	hour := valid.ToString(fmt.Sprintf("%v", time%(3600*24)/3600))
	minute := valid.ToString(fmt.Sprintf("%v", ((time%(3600*24))%3600)/60))
	second := valid.ToString(fmt.Sprintf("%v", ((time%(3600*24))%3600)%60))

	res := ""
	if len(hour) > 0 {
		res += hour + "小时"
	}

	if len(minute) > 0 {
		res += minute + "分"
	}

	if len(second) > 0 {
		res += second + "秒"
	}

	return res
}

/**
格式化时间为local时间
*/
func FormatLocalTime(time time.Time) (str string, err error) {
	jsonTime := models.JsonTime(time)
	if str := strconv.Quote((&jsonTime).String()); len(str) > 0 {
		//去掉引号
		return strconv.Unquote(str)
	}

	return "", errors.New("解析错误")
}

/**
发送任务
*/
func SendDelayQueueRequest(id, order_id, topic_name string) {
	var closeOrder models.CloseOrder
	closeOrder.Topic = topic_name
	closeOrder.ID = id
	closeOrder.Delay = utils.DELAY_JOB_CLOSE
	closeOrder.TTR = utils.DELAY_JOB_TTL
	closeOrder.Body = models.CloseOrderBody{order_id}

	log.Info("the closeOrder is:", closeOrder)
	jsonStr, err := json.Marshal(closeOrder)

	log.Info("the jsonStr is:", jsonStr)

	if err != nil {
		log.Info("the error is:", err.Error())
	}

	url := utils.DELAY_JOB_URL + ":" + utils.DELAY_JOB_PORT + "/" + closeOrder.Topic
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Info("json解析错误:" + err.Error())
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("发送队列任务错误:" + err.Error())
		panic(err)
	}
	defer resp.Body.Close()
}
