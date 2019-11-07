package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func (baseOrm *BaseOrm) AddCartInfo(r *rest.Request) (int, interface{}) {
	course_id, err := valid.ToInt(r.PathParam("course_id"))
	if err != nil {
		log.Info("获取路由参数错误:" + err.Error())
		return 1, "获取路由参数错误:" + err.Error()
	}

	//验证
	var course models.Course
	if err := baseOrm.GetDB().Table("h_edu_courses").Where("id = ? and status = published").First(&course).Error; err != nil {
		log.Info("获取数据错误:" + err.Error())
		return 1, err.Error()
	}

	if course.Type == "free" || course.Price == 0 {
		log.Info("免费课程无法加入购物车")
		return 1, "免费课程无法加入购物车"
	}

	return 1, course_id
}
