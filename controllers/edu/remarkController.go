package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/middlewares"
	log "github.com/sirupsen/logrus"
)

/**
保存评价
 */
type RemarkController struct {
	controller controllers.Controller
}

func (remark *RemarkController) StoreRemark(w rest.ResponseWriter, r *rest.Request) {

	//自定义中间件验证评价内容是否完整，如果完整才会走数据库(类似laravel的Request验证功能)
	var rem middlewares.Remark
	if err := r.DecodeJsonPayload(&rem); err != nil {
		log.Info("参数格式不正确:" + err.Error())
	}

	result, err := (&rem).RemarkValidator()
	if err != nil {
		log.Info("验证错误:" + err.Error())
		remark.controller.JsonReturn(w, "result", err.Error())
	}

	if result {
		//保存评价
		remark.controller.BaseOrm.StoreRemark(r, &rem)
	} else {

	}
}
