package services

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/middlewares"
	"edu_api/models"
)

func (baseOrm *BaseOrm) StoreRemark(r *rest.Request) string {

	var (
		rem         middlewares.Remark
		user_course models.UserCourse
		result      string
	)

	//之前已经验证过了，这里直接使用rem
	r.DecodeJsonPayload(&rem)
	courseId, _ := strconv.Atoi(r.PathParam("id"))

	user = GetUserInfo(r.Header.Get("Authorization"))

	if err := baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", courseId, user.Id).Find(&user_course).Error; err != nil {

		result = "未查到相关课程:" + err.Error()
		log.Info(result)

		return result
	}



	log.Printf("the user id is:%v", user_course.PracticalRating)

	return ""
}
