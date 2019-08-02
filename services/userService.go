package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"strconv"
	"errors"
)

//用户头衔
var userTitle = []string{"暂未选择头衔", "Ph.D", "博士后", "副教授", "教授"}

/**
讲师列表
 */
func (baseOrm *BaseOrm) LecturerList(r *rest.Request) (lectures []models.User, err error) {

	var (
		intLimit int
		id       int
	)
	if limit := r.URL.Query().Get("limit"); len(limit) > 0 {
		intLimit, _ = strconv.Atoi(limit)
	}

	if courseId := r.PathParam("id"); len(courseId) > 0{
		id, _ = strconv.Atoi(courseId)
	}else {
		return lectures, errors.New("课程id必须")
	}

	baseOrm.GetDB().Table("h_users").
		Joins("inner join h_edu_course_user on h_users.id = h_edu_course_user.user_id").
			Where("h_edu_course_user.course_id = ?",id).
				Select("id, avatar, nickname, title").
					Limit(intLimit).
						Find(&lectures)

	for index, _ := range lectures {
		lectures[index].Title = userTitle[0]
	}

	return
}
