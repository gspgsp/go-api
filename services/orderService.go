package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
)

/**
提交订单
*/
func (baseOrm *BaseOrm) SubmitOrder(r *rest.Request, commitOrder *middlewares.CommitOrder) (int, interface{}) {

	var courses []models.Course
	baseOrm.GetDB().Table("h_edu_courses").Where("id in (?)", commitOrder.Ids).Find(&courses)

	if len(courses) == 0 {
		return 1, errors.New("未找到对应ID课程信息")
	}

	return 0, "ok"
}
