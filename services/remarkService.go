package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/utils"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

/**
对于异步请求的返回值 0:成功 1:失败
*/
func (baseOrm *BaseOrm) StoreRemark(r *rest.Request, rem *middlewares.Remark) (int, string) {

	type RatingResult struct {
		AllRating          float64 `json:"all_rating"`
		AllPracticalRating float64 `json:"all_practical_rating"`
		AllPopularRating   float64 `json:"all_popular_rating"`
		AllLogicRating     float64 `json:"all_logic_rating"`
		Number             float64 `json:"number"`
	}

	var (
		user_course models.UserCourse
		rateResult  RatingResult
	)

	courseId, _ := strconv.Atoi(r.PathParam("id"))

	log.Printf("the course_id is:%v", courseId)

	user = GetUserInfo(r.Header.Get("Authorization"))

	baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ? and reviewed = 1", courseId, user.Id).Find(&user_course)

	if user_course.Id > 0 {
		return 1, "该课程已经评价过，不可再次评价"
	}

	//开始评价
	tx := baseOrm.GetDB().Begin()
	sql1 := "update h_user_course set reviewed = 1, anonymous = %d, rating = %s, practical_rating = %d, popular_rating = %d, logic_rating = %d, status = 1, review = %s, reviewed_at = %s where course_id = %d and user_id = %d"

	rating := utils.RetainNumber((rem.PracticalRating + rem.PopularRating + rem.LogicRating) / 3)
	now := models.JsonTime(time.Now())
	reviewed_at := strconv.Quote((&now).String())

	sql1 = fmt.Sprintf(sql1, rem.IsCry, fmt.Sprintf("%.1f", rating), int64(rem.PracticalRating), int64(rem.PopularRating), int64(rem.LogicRating), `"`+rem.Review+`"`, reviewed_at, courseId, user.Id)
	err_1 := tx.Exec(sql1).Error
	//计算课程综合评分
	sql2 := "select sum(practical_rating) as all_practical_rating, sum(popular_rating) as all_popular_rating, sum(logic_rating) as all_logic_rating, sum(rating) as all_rating, count(id) as number from h_user_course where course_id = %d and reviewed = 1 and status = 1 group by course_id"
	sql2 = fmt.Sprintf(sql2, courseId)
	baseOrm.GetDB().Raw(sql2).Scan(&rateResult)

	rateResult.AllLogicRating = utils.RetainNumber((rateResult.AllLogicRating + 10) / (rateResult.Number + 1))
	rateResult.AllPopularRating = utils.RetainNumber((rateResult.AllPopularRating + 10) / (rateResult.Number + 1))
	rateResult.AllPracticalRating = utils.RetainNumber((rateResult.AllPracticalRating + 10) / (rateResult.Number + 1))

	sql3 := "update h_edu_courses set rating = %f, practical_rating = %f, popular_rating = %f, logic_rating = %f where id = %d"
	sql3 = fmt.Sprintf(sql3, rateResult.AllRating, rateResult.AllPracticalRating, rateResult.AllPopularRating, rateResult.AllLogicRating, courseId)
	err_3 := tx.Exec(sql3).Error

	sql4 := "update h_edu_courses set review_count = review_count + 1 where id = %d"
	sql4 = fmt.Sprintf(sql4, courseId)
	err_4 := tx.Exec(sql4).Error

	if err_1 != nil || err_3 != nil || err_4 != nil {
		log.Info("更新课程评分失败，更新用户当前课程评分错误:%s，更新课程评分错误:%s，更新课程评论数错误:%s", err_1.Error(), err_3.Error(), err_4.Error())
		tx.Rollback()
		return 1, "更新课程评分失败"
	} else {
		log.Info("更新课程评分成功")
		tx.Commit()
		return 0, "更新课程评分成功"
	}
}
