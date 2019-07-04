package edu

import (
	"helix-edu-api/services"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"helix-edu-api/models"
)

/**
定义分类控制器
 */
type CategoryController struct {
}

/**
获取所有的分类信息
 */
func (category *CategoryController) GetCategory(w rest.ResponseWriter, r *rest.Request) {
	var (
		err        error
		baseOrm    *services.BaseOrm
		categories []models.Category
	)

	categories, err = baseOrm.CategoryList()

	if err != nil {
		log.Println("query error", err)
	} else {
		returnJson := make(map[string]interface{})

		returnJson["code"] = 0
		returnJson["msg"] = "query successfully!"
		returnJson["categories"] = categories

		w.WriteJson(returnJson)
	}
}
