package edu

import (
	"edu_api/controllers"
	"edu_api/models"
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
	var rollList []models.RollModel

	rollList, exam.controller.Err = exam.controller.BaseOrm.GetExamRollTopicList(r)
	exam.controller.JsonReturn(w, "rollList", rollList)
}
