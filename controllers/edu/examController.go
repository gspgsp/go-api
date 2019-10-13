package edu

import (
	"edu_api/controllers"
	"edu_api/middlewares"
	"edu_api/models"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

/**
作业控制器
*/
type ExamController struct {
	controller controllers.Controller
}

/**
获取题库题目列表
*/
func (exam *ExamController) GetExamRollTopicList(w rest.ResponseWriter, r *rest.Request) {
	var rollList []models.RollModel

	rollList, exam.controller.Err = exam.controller.BaseOrm.GetExamRollTopicList(r)
	exam.controller.JsonReturn(w, "rollList", rollList)
}

/**
获取题库作业详情
*/
func (exam *ExamController) GetExamRollTopicInfo(w rest.ResponseWriter, r *rest.Request) {
	var rollInfo models.RollInfoModel
	rollInfo, exam.controller.Err = exam.controller.BaseOrm.GetExamRollTopicInfo(r)
	exam.controller.JsonReturn(w, "rollInfo", rollInfo)
}

/**
提交答案
*/
func (exam *ExamController) StoreTopicAnswer(w rest.ResponseWriter, r *rest.Request) {
	var ans middlewares.Answer
	if err := r.DecodeJsonPayload(&ans); err != nil {
		log.Info("参数格式不正确:" + err.Error())
	}

	result, err := (&ans).AnswerValidator()
	if err != nil {
		log.Info("验证错误:" + err.Error())
		exam.controller.Err = err
		exam.controller.JsonReturn(w, "result", err.Error())
		return
	}

	if result {
		code, message := exam.controller.BaseOrm.StoreTopicAnswer(r, &ans)
		if code == 0 {
			exam.controller.Err = nil
		} else {
			switch v := message.(type) {
			case string:
				exam.controller.Err = errors.New(v)
			}
		}
		exam.controller.JsonReturn(w, "result", message)
	} else {
		exam.controller.Err = errors.New("未知错误")
		exam.controller.JsonReturn(w, "result", "")
	}
}
