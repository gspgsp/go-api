package controllers

import (
	"edu_api/services"
	"github.com/ant0ine/go-json-rest/rest"
)

/**
控制器下的全局变量
 */
var ReturnJson map[string]interface{}

/**
公共控制器
 */
type Controller struct {
	Err     error
	BaseOrm services.BaseOrm
}

/**
初始化方法
 */
func (that *Controller) init() {
	ReturnJson = make(map[string]interface{})
}

/**
自定义返回json体
 */
func (that *Controller) JsonReturn(w rest.ResponseWriter, baseControl Controller, key interface{}, value interface{}) interface{} {

	//重新初始化
	that.init()

	if baseControl.Err != nil {
		ReturnJson["code"] = 404
		ReturnJson["msg"] = baseControl.Err.Error()
	} else {
		ReturnJson["code"] = 0
		ReturnJson["msg"] = "query successfully!"
		ReturnJson[key.(string)] = value
	}

	w.WriteJson(ReturnJson)

	return nil
}
