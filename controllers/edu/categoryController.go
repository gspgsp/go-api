package edu

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"edu_api/models"
	"edu_api/controllers"
)

/**
定义分类控制器
 */
type CategoryController struct {
	controller controllers.Controller
}

/**
获取所有的分类信息
 */
func (category *CategoryController) GetCategory(w rest.ResponseWriter, r *rest.Request) {
	var (
		categories []models.Category
	)

	categories, category.controller.Err = category.controller.BaseOrm.CategoryList()

	if category.controller.Err != nil {
		log.Println("query error", category.controller.Err)
	} else {
		controllers.ReturnJson["code"] = 0
		controllers.ReturnJson["msg"] = "query successfully!"
		controllers.ReturnJson["categories"] = categories

		w.WriteJson(controllers.ReturnJson)
	}
}
