package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
)

type MaterialController struct {
	controller controllers.Controller
}

/**
资料列表
 */
func (material *MaterialController) GetMaterialList(w rest.ResponseWriter, r *rest.Request) {

	var(
		materials []models.Material
	)

	materials, material.controller.Err = material.controller.BaseOrm.GetMaterialList(r)

	material.controller.JsonReturn(w, "materials", materials)
}
