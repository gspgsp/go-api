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
func (that *Controller) JsonReturn(w rest.ResponseWriter, key interface{}, value interface{}) {

	//重新初始化，因为在main函数的时候，只会调用一次，以后不会再调用，所以限制成当前控制器的方法，重新初始化，否则ReturnJson map里的元素会越来越多
	that.init()

	if that.Err != nil {
		ReturnJson["code"] = 1
		ReturnJson["msg"] = that.Err.Error()
	} else {
		ReturnJson["code"] = 0
		ReturnJson["msg"] = "操作成功!"
		ReturnJson[key.(string)] = value
	}

	w.WriteJson(ReturnJson)
}
