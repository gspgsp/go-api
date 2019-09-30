package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
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
	exam.controller.BaseOrm.GetExamRollTopicList(r)
}
