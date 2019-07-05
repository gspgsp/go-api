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
	Err error
	BaseOrm services.BaseOrm
}

/**
初始化方法
 */
func init()  {
	ReturnJson = make(map[string]interface{})
}