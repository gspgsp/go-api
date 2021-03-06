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
		user = GetUserInfo(r.Header.Get("Authorization"))

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

		now := models.JsonTime(time.Now())
		created_at := strconv.Quote((&now).String())
		updated_at := strconv.Quote((&now).String())
		map_args := map[string]interface{}{"course_id": courseId, "lesson_id": lessonId, "user_id": user.Id, "created_at": created_at}
		if learn.LearnType == 0 {
			if learnCourse.Id > 0 {
				//非视频或音频的时候就+1
				if learn.LesionType == "pdf" || learn.LesionType == "exercise" {
					sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, updated_at, learnCourse.Id)
					updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
				} else {
					sql := fmt.Sprintf("update h_edu_course_learns set updated_at = %s where id = %d", updated_at, learnCourse.Id)
					updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
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
				updateCourseLearn(baseOrm, sql+value, learn.LearnType, map_args)
			}

		} else if learn.LearnType == 1 {
			if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, lesson_length = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, learn.LesionLength, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
			}
		} else if learn.LearnType == 2 {
			if learnCourse.Id > 0 && learnCourse.Status != "finished" {
				sql := fmt.Sprintf("update h_edu_course_learns set status = 2, finish_at = unix_timestamp(now()), watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
			} else if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
			}
		} else if (learn.LearnType == 3) || (learn.LearnType == 4) {
			if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
			}
		} else if learn.LearnType == 4 {
			if learnCourse.Id > 0 {
				sql := fmt.Sprintf("update h_edu_course_learns set watch_num = %d, watch_duration = %d, updated_at = %s where id = %d", learnCourse.WatchNum+1, learnCourse.WatchDuration+learn.WatchDuration, updated_at, learnCourse.Id)
				updateCourseLearn(baseOrm, sql, learn.LearnType, map_args)
			}
		}
	}
}

/**
学习记录更新以及加课程
PC的话，本来learnType等于3(暂停的时候，我是没有更新学习记录以及redis的，因为，页面操作已经包含了所有的离开视频的情况)，但是这里我记录了
 */
func updateCourseLearn(baseOrm *BaseOrm, sql string, learnType int, courseInfo map[string]interface{}) {
	//事务操作
	tx := baseOrm.GetDB().Begin()
	//先看userCourse是否有课，再更新看课记录
	row := baseOrm.GetDB().Table("h_user_course").Where("course_id = ? and user_id = ?", courseInfo["course_id"], courseInfo["user_id"]).Select("id").Row()
	var id int
	row.Scan(&id)

	if id > 0 {
		//待处理，这个地方对课程的唯一可能的操作就是更新一个finished_at(整个课程的完成时间)时间，其它的schedule ... 这些都不处理了，直接放redis
		err_u := tx.Exec(sql).Error

		if err_u != nil {
			log.Info("事务操作出错:" + err_u.Error())
			tx.Rollback()
		} else {
			log.Info("课程已经存在，直接更新")
			tx.Commit()
		}
	} else {
		//插入一条记录，对于免费课和会员免费的课程会出现
		row := baseOrm.GetDB().Table("h_edu_courses").Where("course_id = ?", courseInfo["course_id"]).Select("type").Row()
		var course_type string
		row.Scan(&course_type)

		insert_sql := "insert into `h_edu_courses` (`course_id`, `user_id`, `type`, `created_at`) values"
		insert_value := fmt.Sprintf("(%d,%d,%s,%s)", courseInfo["course_id"], courseInfo["user_id"], course_type, courseInfo["created_at"])

		err_i := tx.Exec(insert_sql + insert_value).Error
		err_u := tx.Exec(sql).Error

		if err_i != nil || err_u != nil {
			log.Info("事务操作出错:" + fmt.Sprintf("插入课程错误:%s,更行课程记录错误:%s", err_i.Error(), err_u.Error()))
			tx.Rollback()
		} else {
			log.Info("课程不存在，先添加再更新")
			tx.Commit()
		}
	}

	//更新学习记录到redis
	defer updateToRedisRecord(baseOrm, courseInfo)
}

/**
更新学习记录到redis
 */
func updateToRedisRecord(baseOrm *BaseOrm, courseInfo map[string]interface{}) {

	row := baseOrm.GetDB().Table("h_edu_courses").Where("id = ?", courseInfo["course_id"]).Select("id, publish_lesson_num").Row()
	var (
		publishLessonNum float64
		courseLearn      []models.CourseLearn
		sumLength        float64
		rate             float64
		watch_duration   int64
	)
	row.Scan(&publishLessonNum)

	err := baseOrm.GetDB().
		Table("h_edu_course_learns").
		Where("user_id = ? and course_id = ? ", courseInfo["user_id"], courseInfo["course_id"]).
		Order("updated_at desc").
		Order("created_at desc").
		Select("id, status, watch_duration").
		Find(&courseLearn).
		Error
	if err != nil {
		log.Info("读取数据错误:" + err.Error())
		return
	}

	if len(courseLearn) <= 0 {
		log.Info("暂无播放记录")
		return
	}

	for _, value := range courseLearn {
		watch_duration += value.WatchDuration
		if value.Status == "finished" {
			sumLength += 1
		} else {
			sumLength += 0.5
		}
	}

	if publishLessonNum > 0 {
		rate = sumLength / publishLessonNum
	}

	SetLatestMediumPlayInfo(courseInfo["user_id"], courseInfo["lesson_id"])

	info := map[string]interface{}{"lesson_id": courseInfo["lesson_id"], "rate": rate, "publish_lesson_num": publishLessonNum, "watch_duration": watch_duration}
	mjson, _ := json.Marshal(info)
	mstring := string(mjson)

	SetClassLatestMediumPlayInfo(courseInfo["user_id"], courseInfo["course_id"], mstring)

	log.Info("缓存观看记录")
}
