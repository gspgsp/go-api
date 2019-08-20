package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"strconv"
	"strings"
	"edu_api/utils"
	"log"
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

	var (
		chapterId = 0
		unitId    = 0
		lessonId  = 0
		where     = make(map[interface{}]interface{}, 5)
	)

	if err := r.DecodeJsonPayload(&learn); err != nil {
		//记录错误日志
	}

	if len(learn.LearnIds) > 0 {
		chapterTypeArray := strings.Split(learn.LearnIds, ":")

		courseId := chapterTypeArray[1]

		if _, err := utils.Contain("a", chapterTypeArray); err == nil {

			chapterId, _ = strconv.Atoi(chapterTypeArray[2])
			unitId, _ = strconv.Atoi(chapterTypeArray[3])
			lessonId, _ = strconv.Atoi(chapterTypeArray[4])

		} else if _, err := utils.Contain("b", chapterTypeArray); err == nil {
			chapterId, _ = strconv.Atoi(chapterTypeArray[2])
			lessonId, _ = strconv.Atoi(chapterTypeArray[3])

		} else if _, err := utils.Contain("c", chapterTypeArray); err == nil {
			unitId, _ = strconv.Atoi(chapterTypeArray[2])
			lessonId, _ = strconv.Atoi(chapterTypeArray[3])

		} else if _, err := utils.Contain("d", chapterTypeArray); err == nil {
			lessonId, _ = strconv.Atoi(chapterTypeArray[2])
		} else if _, err := utils.Contain("s", chapterTypeArray); err == nil {
			lessonId, _ = strconv.Atoi(chapterTypeArray[2])

			if lessonId == 0 {
				row := baseOrm.GetDB().Table("h_edu_chapters").Select("id").Where("course_id = ? and type = 'lesson' and status = 2", courseId).Row()
				row.Scan(&lessonId)
			}

			var parentId = 0
			row := baseOrm.GetDB().Table("h_edu_chapters").Select("parent_id").Where("id = ?", lessonId).Row()
			row.Scan(&parentId)

			if parentId > 0 {
				var chapterType = ""
				row := baseOrm.GetDB().Table("h_edu_chapters").Select("type").Where("id = ?", parentId).Row()
				row.Scan(&chapterType)

				if chapterType == "unit" {

					var preParentId = 0
					row := baseOrm.GetDB().Table("h_edu_chapters").Select("parent_id").Where("id = ?", parentId).Row()
					row.Scan(&preParentId)

					if parentId > 0 {
						chapterId = preParentId
						unitId = parentId
					}

				} else if chapterType == "chapter" {
					chapterId = parentId
				}
			}

		}

		//查询当前视频是否播放过


		where["user_id"] = ""


		log.Printf("the chapterId is:%v\n", chapterId)
		log.Printf("the unitId is:%v\n", unitId)
		log.Printf("the lessonId is:%v\n", lessonId)

		//记录错误日志
		//log.Printf("the course_id is:%v", err.Error())

	}

}
