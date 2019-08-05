package controllers

import "edu_api/services"

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
func init() {
	ReturnJson = make(map[string]interface{})
}

/**
自定义返回json体
 */
func JsonReturn() interface{} {

	/*controllers.ReturnJson = make(map[string]interface{})

	if course.controller.Err != nil {
		controllers.ReturnJson["code"] = 404
		controllers.ReturnJson["msg"] = course.controller.Err.Error()
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["reviews"] = reviews
	}*/




	return nil
}
