package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"strconv"
	"log"
	"strings"
	"edu_api/utils"
)

var (
	num   int64
	learn Learn
)

type Learn struct {
	LearnType     int    `json:"learn_type"`
	LearnIds      string `json:"learn_ids"`
	LesionType    string `json:"lesion_type,omitempty"`
	LesionLength  int    `json:"lesion_length,omitempty"`
	WatchDuration int    `json:"watch_duration,omitempty"`
}

/**
获取播放列表
 */
func (baseOrm *BaseOrm) GetPlayList(r *rest.Request) (playList models.Media, err error) {

	id, err := strconv.Atoi(r.PathParam("id"))
	lesionId, err := strconv.Atoi(r.PathParam("lesion_id"))
	if err != nil {
		return
	}

	chapterList := make([]models.Chapter, 10)

	first := models.Chapter{}

	if err := baseOrm.GetDB().Table("h_edu_chapters").Where("course_id = ? and status = 'published'", id).Find(&chapterList).Error; err != nil || len(chapterList) == 0 {
		return playList, err
	}

	//当前课程类型
	type courseType struct {
		Type string `json:"type"`
	}
	cType := courseType{}

	baseOrm.GetDB().Table("h_edu_chapters").Select("h_edu_courses.type").Joins("left join h_edu_courses on h_edu_courses.id = h_edu_chapters.course_id").Where("h_edu_chapters.id = ?", lesionId).Find(&cType)

	if err := baseOrm.GetDB().Table("h_edu_chapters").Where("id = ? and status = 'published'", lesionId).Find(&first).Error; err != nil {
		return playList, err
	}

	if first.Id == 0 {
		//就选择当前课程下最远的一条数据
		if err := baseOrm.GetDB().Table("h_edu_chapters").Where("course_id = ? and status = 'published' and type = 'lesson'", id).Order("updated_at asc").Order("created_at asc").Find(&first).Error; err != nil || first.Id == 0 {
			return playList, err
		}
	}

	//数据组合
	playList.Chapter = Trees(chapterList).([]models.Chapter)

	if cType.Type == "class" {
		//获取当前lesion所在的章节
		row := baseOrm.GetDB().Table("h_edu_chapters").Select("title").Where("id = ?", recursion(baseOrm, lesionId)).Row()
		row.Scan(&playList.CurrentTitle)
	}

	//根据类型判断
	playList.CurrentLesion = first.Title
	playList.LesionType = cType.Type

	return
}

/**
递归取父级
 */
func recursion(baseOrm *BaseOrm, parent_id int) (id int) {

	row := baseOrm.GetDB().Raw("select parent_id from h_edu_chapters where id =?", parent_id).Row()
	row.Scan(&id)

	num++

	if num == 3 {
		return id
	} else {
		return recursion(baseOrm, id)
	}
}

/**
a:三级 b/c:二级 d:一级 s:自定义，未点击事件
 */
func (baseOrm *BaseOrm) PutCourseLearn(r *rest.Request) {

	if err := r.DecodeJsonPayload(&learn); err != nil {
		//记录错误日志
	}

	if len(learn.LearnIds) > 0 {
		chapter_type_array := strings.Split(learn.LearnIds, ":")

		course_id := chapter_type_array[1]

		if _, err := utils.Contain("s", chapter_type_array); err == nil {
			log.Printf("the course_id is:%v", course_id)
		} else {
			log.Printf("the course_id is:%v", "afadsf")
		}
	}

}
