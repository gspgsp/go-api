package services

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
	"helix-edu-api/utils"
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

	err = jsonStruct.Load("E:/gorespority/src/helix-edu-api/config/database.json", &v)

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
