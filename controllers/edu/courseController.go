package edu

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/services"
	"edu_api/models"
	"log"
)

type CourseController struct {
}

func (course *CourseController) GetCourseList(w rest.ResponseWriter, r *rest.Request) {
	var (
		err     error
		baseOrm *services.BaseOrm
		courses []models.Course
	)

	courses, err = baseOrm.CourseList(r)

	if err != nil {
		log.Println("query error", err)
	} else {
		returnJson := make(map[string]interface{})

		returnJson["code"] = 0
		returnJson["msg"] = "query successfully!"
		returnJson["courses"] = courses

		w.WriteJson(returnJson)
	}

}

func (course *CourseController) GetRecommendList(w rest.ResponseWriter, r *rest.Request) {

}
