package edu

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
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

	course.controller.JsonReturn(w, "courses", courses)
}

/**
获取课程详情列表
 */
func (course *CourseController) GetCourseDetail(w rest.ResponseWriter, r *rest.Request) {
	var (
		detail models.Detail
	)

	detail, course.controller.Err = course.controller.BaseOrm.GetCourseDetail(r)

	course.controller.JsonReturn(w, "detail", detail)
}

/**
获取套餐列表
 */
func (course *CourseController) GetPackageList(w rest.ResponseWriter, r *rest.Request) {
	var (
		packages []models.Package
	)

	packages, course.controller.Err = course.controller.BaseOrm.PackageList(r)

	course.controller.JsonReturn(w, "packages", packages)
}

/**
获取课程章节
 */
func (course *CourseController) GetCourseChapter(w rest.ResponseWriter, r *rest.Request) {
	var (
		chapters []models.Chapter
	)

	chapters, course.controller.Err = course.controller.BaseOrm.GetCourseChapter(r)

	course.controller.JsonReturn(w, "chapters", chapters)
}

/**
课程评价列表
 */
func (course *CourseController) GetCourseReview(w rest.ResponseWriter, r *rest.Request) {
	var reviews []models.Review

	reviews, course.controller.Err = course.controller.BaseOrm.GetCourseReview(r)

	course.controller.JsonReturn(w, "reviews", reviews)
}

/**
推荐课程
 */
func (course *CourseController) GetRecommendCourse(w rest.ResponseWriter, r *rest.Request) {
	var recommends []models.Recommend

	recommends, course.controller.Err = course.controller.BaseOrm.GetRecommendCourse(r)

	course.controller.JsonReturn(w, "recommends", recommends)
}
