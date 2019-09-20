package services

import (
	"strconv"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/middlewares"
	"edu_api/models"
	"fmt"
	"edu_api/utils"
	"time"
	log "github.com/sirupsen/logrus"
)

func (baseOrm *BaseOrm) StoreRemark(r *rest.Request, rem *middlewares.Remark) string {

	var (
		user_course models.UserCourse
	)

	courseId, _ := strconv.Atoi(r.PathParam("id"))

	log.Printf("the course_id is:%v", courseId)

	user = GetUserInfo(r.Header.Get("Authorization"))

	baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ? and reviewed = 1", courseId, user.Id).Find(&user_course)

	if user_course.Id > 0 {
		return "该课程已经评价过，不可再次评价"
	}

	//开始评价
	//tx := baseOrm.GetDB().Begin()
	sql1 := "update h_user_course set reviewed = 1, anonymous = %d, rating = %s, practical_rating = %d, popular_rating = %d, logic_rating = %d, status = 1, review = %s, reviewed_at = %s where course_id = %d and user_id = %d"

	rating := utils.RetainNumber((rem.PracticalRating + rem.PopularRating + rem.LogicRating) / 3)
	now := models.JsonTime(time.Now())
	reviewed_at := strconv.Quote((&now).String())

	sql1 = fmt.Sprintf(sql1, rem.IsCry, fmt.Sprintf("%.1f", rating), int64(rem.PracticalRating), int64(rem.PopularRating), int64(rem.LogicRating), rem.Review, reviewed_at, courseId, user.Id)

	log.Printf("the sql is:%s", sql1)
	return ""
}
