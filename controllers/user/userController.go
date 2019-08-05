package user

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"log"
)

type UserController struct {
	controller controllers.Controller
}

/**
讲师列表
 */
func (user *UserController) GetLecturerList(w rest.ResponseWriter, r *rest.Request) {
	var (
		lecturers []models.User
	)

	lecturers, user.controller.Err = user.controller.BaseOrm.LecturerList(r)

	controllers.ReturnJson = make(map[string]interface{})
	if user.controller.Err != nil {
		log.Println("query error", user.controller.Err.Error())
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["lecturers"] = lecturers

		w.WriteJson(controllers.ReturnJson)
	}
}
