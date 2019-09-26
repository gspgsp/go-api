package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
)

type PackageController struct {
	controller controllers.Controller
}

func (course *PackageController) GetComposePackage(w rest.ResponseWriter, r *rest.Request) {
	course.controller.BaseOrm.GetComposePackage(r)
}
