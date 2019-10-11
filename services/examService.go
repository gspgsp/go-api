package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"encoding/json"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
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
func (baseOrm *BaseOrm) StoreTopicAnswer(r *rest.Request, answer *middlewares.Answer) (int, string) {
	var (
		userCourse models.UserCourse
		topics     []models.TopicModel
		topicIds   []int64
	)

	user = GetUserInfo(r.Header.Get("Authorization"))
	baseOrm.GetDB().Table("h_user_course").Where("user_id = ? and course_id = ?", user.Id, answer.CourseId).Select("id").First(&userCourse)
	if userCourse.Id > 0 {
		for _, value := range answer.Answers {
			topicIds = append(topicIds, value.TopicId)
		}
		//如果为0，则一定为自动提交(规定时间内没答题目)；当不为0的时候，可能为自动提交(没在规定时间内答完题目)
		if len(topicIds) == 0 {
			//中间表查询
		} else {
			//if err1 := baseOrm.GetDB().Table("h_exam_topics").Where("id in (?)", topicIds).Select("id, options").Find(&topics).Error; err1 != nil {
			//	log.Info("查询错误:" + err1.Error())
			//	return 1, "查询错误:" + err1.Error()
			//}

			if err1 := baseOrm.GetDB().Table("h_exam_topics").Joins("left join h_exam_roll_topic on h_exam_roll_topic.topic_id = h_exam_topics.id").Where("h_exam_topics.status = 1 and h_exam_roll_topic.roll_id = ?", answer.RollId).Select("h_exam_topics.id, h_exam_topics.options").Find(&topics).Error; err1 != nil {
				log.Info("查询错误:" + err1.Error())
				return 1, "查询错误:" + err1.Error()
			}

			//异步查询
			var channel = make(chan []models.GradeLogResult, 4)
			var endChannel = make(chan bool)
			endNumber := 0
			number := len(topics)
			var wg sync.WaitGroup
			for i := 0; i < len(topics); i++ {
				wg.Add(1)
				go judgeAnswerResult(channel, endChannel, baseOrm, topics[i].Id, topics[i].Options, answer.Answers, user.Id, answer.CourseId)
			}
			wg.Wait()
		}

		//label

		return 0, "操作成功"
	} else {
		return 1, "当前用户没有此课程"
	}
}

func judgeAnswerResult(channel chan []models.GradeLogResult, endChannel chan bool, baseOrm *BaseOrm, topicId int64, options string, userOption []middlewares.AnswerData, userId int, courseId int64) {
	defer func() {
		err := recover()
		if err != nil {
			log.Info("运行异常:")
		}
		endChannel <- true
	}()

	var parseOptions []models.OptionModel
	for _, value := range userOption {
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

			if value.Option == strings.Join(rightSli, "|") {

			} else {

			}

		}
	}
}
