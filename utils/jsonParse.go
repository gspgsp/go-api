package utils

import (
	"io/ioutil"
	"encoding/json"
)

type JsonStruct struct {
}

/**
格式化json
 */
func InitJsonStruct() *JsonStruct {

	return &JsonStruct{}
}

/**
解析json配置文件到指定对象，这里用户interface标识
 */
func (that *JsonStruct) Load(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	dataJson := []byte(data)

	if err = json.Unmarshal(dataJson, v); err != nil {
		return err
	}

	return  nil
}

/**
获取mysql数据库配置
 */
type ConfigType struct {
	Mysql mysqlConfigValue `json:"mysql"`
}

/**
mysql护具库配置详情
 */
type mysqlConfigValue struct {
	DbConnect  string `json:"db_connect"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbDatabase string `json:"db_database"`
	DbUsername string `json:"db_username"`
	DbPassword string `json:"db_password"`
	DbPrefix   string `json:"db_prefix"`
}
