package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"strconv"
	"strings"
	"edu_api/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
	"fmt"
)

var (
	num   int64
	learn Learn
	//user models.User
)

type Learn struct {
	LearnType     int    `json:"learn_type"`
	LearnIds      string `json:"learn_ids"`
	LesionType    string `json:"lesion_type,omitempty"`
	LesionLength  int    `json:"lesion_length,omitempty"`
	WatchDuration int64  `json:"watch_duration,omitempty"`
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
		where     = make(map[string]interface{}, 5) //key一定要是string类型，否则where条件查询会出错
	)

	if err := r.DecodeJsonPayload(&learn); err != nil {
		//记录错误日志
	}

	if len(learn.LearnIds) > 0 {
		chapterTypeArray := strings.Split(learn.LearnIds, ":")

		courseId, _ := strconv.Atoi(chapterTypeArray[1])

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
		authHeaderToken := r.Header.Get("Authorization")
		info := GetRedisCache(authHeaderToken, "hget", "info")

		//这个user变量在登录的时候已经声明过了
		if err := json.Unmarshal([]byte(info), &user); err != nil {
			//记录日志
			log.Info("解析用户信息错误:", err.Error())
		}
		where["user_id"] = user.Id
		where["course_id"] = courseId
		where["lesson_id"] = lessonId

		if chapterId > 0 {
			where["chapter_id"] = chapterId
		}

		if unitId > 0 {
			where["unit_id"] = unitId
		}

		var learnCourse models.CourseLearn
		baseOrm.GetDB().Table("h_edu_course_learns").Where(where).First(&learnCourse)

		//b,_:=json.Marshal(learnCourse)，将learnCourse转为json格式输出，其中的time要格式化为本地时间格式（否则输出的时间格式是西方时间格式），采用的方式是 重写MarshalJSON，
		now := models.JsonTime(time.Now())
		created_at := strconv.Quote((&now).String())
		updated_at := strconv.Quote((&now).String())
		if learn.LearnType == 0 {
			if learnCourse.Id > 0 {
				//非视频或音频的时候就+1
				if learn.LesionType == "pdf" || learn.LesionType == "exercise" {
					sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, updated_at, learnCourse.Id)
					updateCourseLearn(baseOrm, sql)
				} else {
					sql := fmt.Sprintf("update h_edu_course_learns set updated_at = %s where id = %d", updated_at, learnCourse.Id)
					updateCourseLearn(baseOrm, sql)
				}
			} else {
				//新增一条数据
				sql := "insert into `h_edu_course_learns` (`status`, `start_at`,`finish_at`,`watch_duration`,`lesson_length`,`watch_num`,`user_id`,`course_id`,`chapter_id`,`unit_id`,`lesson_id`,`created_at`,`updated_at`) values "
				watch_num := 0
				status := 1
				finish_at := "NULL"

				if learn.LesionType == "pdf" || learn.LesionType == "exercise" {
					watch_num = 1
					status = 2
					finish_at = strconv.FormatInt(time.Now().Unix(), 10)
				}

				value := fmt.Sprintf("(%d,%d,%s,%d,%d,%d,%d,%d,%d,%d,%d,%s,%s)", status, time.Now().Unix(), finish_at, 0, learn.LesionLength, watch_num, user.Id, courseId, chapterId, unitId, lessonId, created_at, updated_at)
				updateCourseLearn(baseOrm, sql + value)
			}

		} else if learn.LearnType == 1 {
			if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, lesson_length = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, learn.LesionLength, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql)
			}
		} else if learn.LearnType == 2 {
			if learnCourse.Id > 0 && learnCourse.Status != "finished" {
				sql := fmt.Sprintf("update h_edu_course_learns set status = 2, finish_at = unix_timestamp(now()), watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql)
			} else if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql)
			}
		} else if (learn.LearnType == 3) || (learn.LearnType == 4) {
			if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql)
			}
		} else if learn.LearnType == 4 {
			if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql)
			}
		}

		//更新学习记录到redis
		defer updateToRedisRecord()
	}
}

func updateCourseLearn(baseOrm *BaseOrm, sql string) {
	tx := baseOrm.GetDB().Begin()
	//先看userCourse是否有课，再更新看课记录
	err := tx.Exec(sql).Error

	if err != nil {
		tx.Rollback()
	} else {
		log.Info("可以通过事务提交")
		tx.Commit()
	}
}

func updateToRedisRecord() {

	log.Info("到defer了......")
}
