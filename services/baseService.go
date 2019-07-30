package services

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
	"edu_api/utils"
	"edu_api/models"
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

	err = jsonStruct.Load("E:/gorespority/src/edu_api/config/database.json", &v)

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
	that.DB.DB().SetMaxIdleConns(10)
	that.DB.DB().SetMaxOpenConns(100)
	that.DB.DB().SetConnMaxLifetime(300 * time.Second)
	that.DB.LogMode(true)

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
				//用当前id去取值，parent_id会和id对应起来
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
