package user

import (
	"edu_api/controllers"
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
)

type UserController struct {
	controller controllers.Controller
}

/**
用户课程讲师列表
*/
func (user *UserController) GetLecturerList(w rest.ResponseWriter, r *rest.Request) {
	var (
		lecturers []models.User
	)

	lecturers, user.controller.Err = user.controller.BaseOrm.LecturerList(r)

	user.controller.JsonReturn(w, "lecturers", lecturers)
}
