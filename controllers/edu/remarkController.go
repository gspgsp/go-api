package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
)

/**
保存评价
 */
type RemarkController struct {
	controller controllers.Controller
}

func (remark *RemarkController) StoreRemark(w rest.ResponseWriter, r *rest.Request) {

	//自定义中间件验证评价内容是否完整，如果完整才会走数据库(类似laravel的Request验证功能)


	//保存评价
	remark.controller.BaseOrm.StoreRemark(r)
}
