package user

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
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

	user.controller.JsonReturn(w, "lecturers", lecturers)
}
