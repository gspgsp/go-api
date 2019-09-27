package edu

import (
	"edu_api/controllers"
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
)

type PackageController struct {
	controller controllers.Controller
}

/**
组合套餐
*/
func (course *PackageController) GetComposePackage(w rest.ResponseWriter, r *rest.Request) {
	var (
		compose models.ComposeModel
	)
	compose, course.controller.Err = course.controller.BaseOrm.GetComposePackage(r)
	course.controller.JsonReturn(w, "compose", compose.ComposePackage)
}
