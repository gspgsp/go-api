package edu

import "github.com/ant0ine/go-json-rest/rest"

type CourseController struct {
}

type Filter struct {
	Code            string
	DifficultyLevel int
	Order           string
	VipPrice        int
}

func (course *CourseController) GetCourseList(w rest.ResponseWriter, r *rest.Request) {

}
