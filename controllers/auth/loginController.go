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

	if login.controller.Err != nil {
		controllers.ReturnJson["code"] = 404
		controllers.ReturnJson["msg"] = login.controller.Err
	}else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["token"] = token
	}

	w.WriteJson(controllers.ReturnJson)
}
