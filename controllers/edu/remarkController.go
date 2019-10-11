package edu

import (
	"edu_api/controllers"
	"edu_api/middlewares"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

/**
保存评价
*/
type RemarkController struct {
	controller controllers.Controller
}

/**
创建评价
*/
func (remark *RemarkController) StoreRemark(w rest.ResponseWriter, r *rest.Request) {

	//自定义中间件验证评价内容是否完整，如果完整才会走数据库(类似laravel的Request验证功能)
	var rem middlewares.Remark
	if err := r.DecodeJsonPayload(&rem); err != nil {
		log.Info("参数格式不正确:" + err.Error())
	}

	result, err := (&rem).RemarkValidator()
	if err != nil {
		log.Info("验证错误:" + err.Error())
		remark.controller.Err = err
		remark.controller.JsonReturn(w, "result", err.Error())
		return
	}

	if result {
		//保存评价
		code, message := remark.controller.BaseOrm.StoreRemark(r, &rem)
		if code == 0 {
			remark.controller.Err = nil
		} else {
			remark.controller.Err = errors.New(message)
		}

		remark.controller.JsonReturn(w, "result", message)
	} else {
		remark.controller.Err = errors.New("未知错误")
		remark.controller.JsonReturn(w, "result", "")
	}
}
