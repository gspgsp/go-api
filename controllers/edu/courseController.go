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

	controllers.ReturnJson = make(map[string]interface{})
	if course.controller.Err != nil {
		log.Println("query error", course.controller.Err)
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["courses"] = courses

		w.WriteJson(controllers.ReturnJson)
	}
}

/**
获取课程详情列表
 */
func (course *CourseController) GetCourseDetail(w rest.ResponseWriter, r *rest.Request) {
	var (
		detail models.Detail
	)

	detail, course.controller.Err = course.controller.BaseOrm.GetCourseDetail(r)

	controllers.ReturnJson = make(map[string]interface{})
	if course.controller.Err != nil {
		controllers.ReturnJson["code"] = 404
		controllers.ReturnJson["msg"] = course.controller.Err.Error()
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["detail"] = detail
	}

	w.WriteJson(controllers.ReturnJson)
}

/**
获取套餐列表
 */
func (course *CourseController) GetPackageList(w rest.ResponseWriter, r *rest.Request) {
	var (
		packages []models.Package
	)

	packages, course.controller.Err = course.controller.BaseOrm.PackageList(r)

	controllers.ReturnJson = make(map[string]interface{})
	if course.controller.Err != nil {
		controllers.ReturnJson["code"] = 404
		controllers.ReturnJson["msg"] = course.controller.Err.Error()
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["courses"] = packages
	}

	w.WriteJson(controllers.ReturnJson)
}

/**
获取课程章节
 */
func (course *CourseController) GetCourseChapter(w rest.ResponseWriter, r *rest.Request) {
	var (
		chapters []models.Chapter
	)

	chapters, course.controller.Err = course.controller.BaseOrm.GetCourseChapter(r)

	controllers.ReturnJson = make(map[string]interface{})
	if course.controller.Err != nil {
		controllers.ReturnJson["code"] = 404
		controllers.ReturnJson["msg"] = course.controller.Err.Error()
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["chapters"] = chapters
	}

	w.WriteJson(controllers.ReturnJson)
}
