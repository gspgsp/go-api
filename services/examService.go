package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
	"sort"
)

/**
获取题库题目列表
*/
func (baseOrm *BaseOrm) GetExamRollTopicList(r *rest.Request) (rollList []models.RollModel, err error) {
	var (
		gradeChan  chan models.GradeModel
		endChan    chan bool
		endChanNum int
		allChanNum int
	)

	id, err1 := valid.ToInt(r.PathParam("id"))
	if err1 != nil {
		log.Info("获取请求参数错误:" + err1.Error())
		return nil, err1
	}

	if err2 := baseOrm.GetDB().Table("h_exam_rolls").Where("course_id = ? and status = 2", id).Find(&rollList).Error; err != nil {
		log.Info("获取数据错误:" + err2.Error())
		return nil, err2
	}

	user = GetUserInfo(r.Header.Get("Authorization"))
	//计算是否答过当前考卷，以及考试结果
	if len(rollList) == 0 {
		return nil, errors.New("当前课程无考卷信息")
	}

	allChanNum = len(rollList)
	gradeChan = make(chan models.GradeModel, 4)
	endChan = make(chan bool)
	for _, value := range rollList {
		go getGrade(baseOrm, value.Id, user.Id, gradeChan, endChan)
	}

	/**
	  label处理channel
	*/
L:
	for {
		select {
		case grade, ok := <-gradeChan:
			if ok {
				for index, value := range rollList {
					if value.Id == grade.RollId {
						rollList[index].Grade = grade
					}
				}
			}
		case <-endChan:
			endChanNum++
			if endChanNum == allChanNum {
				close(gradeChan)
				close(endChan)
				break L
			}
		}
	}

	//数据单独处理
	for i, val := range rollList {
		if val.Mode == 1 { //随机模式，计算总分
			score := 0
			allScore := 0

			rows, err := baseOrm.GetDB().
				Table("h_exam_topics").
				Joins("left join h_exam_roll_topic on h_exam_topics.id = h_exam_roll_topic.topic_id").
				Where("h_exam_roll_topic.roll_id = ?", val.Id).Select("h_exam_topics.score").Rows()
			if err == nil {
				for rows.Next() {
					rows.Scan(&score)
					allScore += score
				}
			}
			rollList[i].TotalScore = allScore
		} else if val.Mode == 2 { //分值模式，计算总数
			//暂无数据
		}

		//未答过，需要格式化答题时间为分钟；答过的题目直接给结果grade里面有结果
		if val.Grade.Id == 0 {
			rollList[i].LimitedAt, _ = FormatTime(val.LimitedAt)
		}
	}

	return rollList, nil
}

/**
获取当前题库成绩
*/
func getGrade(baseOrm *BaseOrm, rollId int64, userId int, gradeChan chan models.GradeModel, endChan chan bool) {
	var grade models.GradeModel

	defer func() {
		endChan <- true
	}()

	baseOrm.GetDB().Table("h_exam_grades").Where("roll_id = ? and user_id = ?", rollId, userId).First(&grade)
	//grade.Id = 1
	//grade.RollId = 2
	//grade.CourseId = 2
	//grade.ChapterId = 3
	//grade.Point = 4
	//grade.Result = `{"point": 4, "numbers": 10, "success": 4, "usetimes": 16, "all_point": 10}`
	gradeChan <- grade
}

/**
获取题库作业详情
*/
func (baseOrm *BaseOrm) GetExamRollTopicInfo(r *rest.Request) (rollInfo models.RollInfoModel, err error) {

	var (
		topics []models.TopicModel
	)
	id, err1 := valid.ToInt(r.PathParam("id"))
	if err1 != nil {
		log.Info("获取请求参数错误:" + err1.Error())
		return rollInfo, err1
	}

	if err2 := baseOrm.GetDB().Table("h_exam_rolls").Where("id = ? and status = 2", id).First(&rollInfo).Error; err2 != nil {
		log.Info("获取数据错误:" + err2.Error())
		return rollInfo, err2
	}

	if err3 := baseOrm.GetDB().
		Table("h_exam_topics").
		Joins("left join h_exam_roll_topic on h_exam_topics.id = h_exam_roll_topic.topic_id").
		Where("h_exam_roll_topic.roll_id = ?", id).Select("h_exam_topics.*").Find(&topics).Error; err3 != nil {
		log.Info("获取数据错误:" + err3.Error())
		return rollInfo, err3
	}

	var v []models.OptionModel
	for i, val := range topics {
		if err4 := json.Unmarshal([]byte(val.Options), &v); err4 != nil {
			log.Info("解析答案错误:" + err4.Error())
			return rollInfo, err4
		}

		for _, value := range v {
			topics[i].ParseOptions = append(topics[i].ParseOptions, value)
		}

		//手动忽略掉json的option
		topics[i].Options = ""
	}

	rollInfo.Topics = topics

	return rollInfo, nil
}

/**
提交答案
*/
func (baseOrm *BaseOrm) StoreTopicAnswer(r *rest.Request, answer *middlewares.Answer) (int, interface{}) {
	var (
		userCourse     models.UserCourse
		grade          models.GradeModel
		topics         []models.TopicModel
		gradeLogResult []models.GradeLogResult
		allPoint       int64
		point          int64
		success        int
		numbers        int
	)

	user = GetUserInfo(r.Header.Get("Authorization"))
	baseOrm.GetDB().Table("h_user_course").Where("user_id = ? and course_id = ?", user.Id, answer.CourseId).Select("id").First(&userCourse)
	if userCourse.Id > 0 {
		//是否答过当前试卷
		baseOrm.GetDB().Table("h_exam_grades").Where("roll_id = ? and course_id = ? and user_id = ? ", answer.RollId, answer.CourseId, user.Id).Select("id").First(&grade)
		if grade.Id > 0 {
			log.Info("当前试卷已答过，请勿重复答题")
			return 1, "当前试卷已答过，请勿重复答题"
		}

		//是否存在题目
		if err1 := baseOrm.GetDB().Table("h_exam_topics").Joins("left join h_exam_roll_topic on h_exam_roll_topic.topic_id = h_exam_topics.id").Where("h_exam_topics.status = 1 and h_exam_roll_topic.roll_id = ?", answer.RollId).Select("h_exam_topics.id, h_exam_topics.options, h_exam_topics.score").Find(&topics).Error; err1 != nil {
			log.Info("查询错误:" + err1.Error())
			return 1, "查询错误:" + err1.Error()
		}

		//异步查询
		var channel = make(chan models.GradeLogResult, 30)
		var endChannel = make(chan bool)
		endNumber := 0
		numbers = len(topics)
		for i := 0; i < len(topics); i++ {
			allPoint += topics[i].Score
			go judgeAnswerResult(channel, endChannel, baseOrm, topics[i].Id, topics[i].Options, answer.Answers, user.Id, answer.CourseId)
		}

		//label
	GradeLogResultChannelLabel:
		for {
			select {
			case v, ok := <-channel:
				if ok {
					gradeLogResult = append(gradeLogResult, v)
				}
			case <-endChannel:
				endNumber++
				if endNumber == numbers {
					close(channel)
					break GradeLogResultChannelLabel
				}
			}
		}
		//统计结果
		for _, value := range gradeLogResult {
			if value.Num == value.UserChose {
				success += 1
				for _, val := range topics {
					if val.Id == value.TopicId {
						point += val.Score
					}
				}
			}
		}

		//插入数据
		insert_sql := "insert into `h_exam_grades` (`point`, `result`, `created_at`, `updated_at`, `roll_id`, `course_id`, `user_id`) values"
		gradeResult := models.GradeResult{point, numbers, success, answer.StartTime, allPoint}
		jsonStr, _ := json.Marshal(gradeResult)

		now := models.JsonTime(time.Now())
		created_at := strconv.Quote((&now).String())
		updated_at := strconv.Quote((&now).String())
		fmt.Printf("the string(jsonStr) is:%v", "`"+string(jsonStr)+"`")
		insert_value := fmt.Sprintf("(%d,%s,%s,%s,%d,%d,%d)", point, "'"+string(jsonStr)+"'", created_at, updated_at, answer.RollId, answer.CourseId, user.Id)

		tx := baseOrm.GetDB().Begin()
		//因为gorm将golang自带的database/sql库的 Lastinsertid方法封装掉了，本来执行exec方法以后，可以返回结果集以取到最新一条数据的id;但是gorm操作下却不行
		//所以我想到下面的链式操作，通过Last方法获取到最新一条插入的数据，这个因为是在tx下的操作，应该不会有问题
		var grade models.GradeModel
		err1 := tx.Exec(insert_sql + insert_value).Table("h_exam_grades").Last(&grade).Error

		var err2 error
		var is_correct = 0
		var gradeLogResultSlice models.GradeLogResultSlice = gradeLogResult

		//重新排序，goroutine导致index混乱
		var answerResultReturn []models.AnswerResultReturn
		sort.Stable(gradeLogResultSlice)
		for i := 0; i < len(gradeLogResultSlice); i++ {
			if gradeLogResultSlice[i].Num == gradeLogResultSlice[i].UserChose {
				is_correct = 1
			} else {
				is_correct = 0
			}

			//最终返回结果集
			answerResultReturn = append(answerResultReturn, models.AnswerResultReturn{(i + 1), is_correct})

			gradeLog := models.GradeLogResult{Num: gradeLogResultSlice[i].Num, UserChose: gradeLogResultSlice[i].UserChose}
			jsonLogStr, _ := json.Marshal(gradeLog)
			insert_grade_sql := "insert into `h_exam_grade_logs` (`is_correct`, `result`, `created_at`, `updated_at`, `grade_id`, `roll_id`, `course_id`, `topic_id`, `user_id`) values"
			value := fmt.Sprintf("(%d,%s,%s,%s,%d,%d,%d,%d,%d)", is_correct, "'"+string(jsonLogStr)+"'", created_at, updated_at, grade.Id, answer.RollId, answer.CourseId, gradeLogResultSlice[i].TopicId, user.Id)
			err2 = tx.Exec(insert_grade_sql + value).Error
		}

		if err1 != nil || err2 != nil {
			log.Info("事务操作出错:" + fmt.Sprintf("插入答题记录错误:%s", err1.Error()))
			tx.Rollback()
		} else {
			log.Info("插入答题记录成功")
			tx.Commit()
		}

		//返回结果
		var answerReturn models.AnswerReturn
		answerReturn.Point = point
		answerReturn.AllPoint = allPoint
		answerReturn.Success = success
		answerReturn.Numbers = numbers
		answerReturn.SubmitTime = created_at
		//answerReturn.UseTime = FormatTimeToChinese(time.Now().Unix() - answer.StartTime)
		answerReturn.Result = answerResultReturn
		return 0, answerReturn
	} else {
		return 1, "当前用户没有此课程"
	}
}

func judgeAnswerResult(channel chan models.GradeLogResult, endChannel chan bool, baseOrm *BaseOrm, topicId int64, options string, userOption []middlewares.AnswerData, userId int, courseId int64) {
	defer func() {
		err := recover()
		if err != nil {
			log.Info("goroutine运行异常:")
		}
		endChannel <- true
	}()

	var parseOptions []models.OptionModel
	for _, value := range userOption {
		//这里要求用户答的题目数和试卷上的一样，那么就一定可以找到对应的结果，就不用了处理else情况了
		if value.TopicId == topicId {
			err := json.Unmarshal([]byte(options), &parseOptions)
			if err != nil {
				log.Info("题目选项解析错误:" + err.Error())
				return
			}

			var rightSli []string
			for _, val := range parseOptions {
				if val.IsRight == "1" {
					rightSli = append(rightSli, val.Num)
				}
			}
			//记录日志
			var gradeLogResult = models.GradeLogResult{topicId, strings.Join(rightSli, ","), value.Option}
			channel <- gradeLogResult
		}
	}
}
