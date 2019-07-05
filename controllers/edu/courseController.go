package edu

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"log"
	"edu_api/controllers"
)

/**
定义课程控制器
 */
type CourseController struct {
	controller controllers.Controller
}

/**
获取所有课程信息
 */
func (course *CourseController) GetCourseList(w rest.ResponseWriter, r *rest.Request) {
	var (
		courses []models.Course
	)

	courses, course.controller.Err = course.controller.BaseOrm.CourseList(r)

	if course.controller.Err != nil {
		log.Println("query error", course.controller.Err)
	} else {
		returnJson := make(map[string]interface{})

		returnJson["code"] = 0
		returnJson["msg"] = "query successfully!"
		returnJson["courses"] = courses

		w.WriteJson(returnJson)
	}

}

func (course *CourseController) GetPackageList(w rest.ResponseWriter, r *rest.Request) {

}
