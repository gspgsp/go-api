package auth

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
)

type LoginController struct {
	controller controllers.Controller
}

func (login *LoginController) Login(w rest.ResponseWriter, r *rest.Request) {

	var (
		token string
	)

	token, login.controller.Err = login.controller.BaseOrm.Login(r)

	login.controller.JsonReturn(w, login.controller, "token", token)
}
