package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func (baseOrm *BaseOrm) GetExamRollTopicList(r *rest.Request) {
	var (
		rollList []models.RollModel
	)

	id, err := valid.ToInt(r.PathParam("id"))
	if err != nil {
		return
	}

	if err := baseOrm.GetDB().Table("h_exam_rolls").Where("course_id = ? and status = 2", id).Find(&rollList).Error; err != nil {
		log.Info("获取数据错误:" + err.Error())
		return
	}

	user = GetUserInfo(r.Header.Get("Authorization"))

	log.Printf("the user info is:%v", user)
}
