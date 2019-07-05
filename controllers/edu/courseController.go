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
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["courses"] = courses

		w.WriteJson(controllers.ReturnJson)
	}

}

func (course *CourseController) GetPackageList(w rest.ResponseWriter, r *rest.Request) {
	var (
		packages []models.Package
	)

	packages, course.controller.Err = course.controller.BaseOrm.PackageList()

	if course.controller.Err != nil {
		log.Println("query error", course.controller.Err)
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["courses"] = packages

		w.WriteJson(controllers.ReturnJson)
	}

}
