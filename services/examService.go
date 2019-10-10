package services

import (
	"edu_api/models"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
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

	rollInfo.Topics = topics
	log.Printf("the rollInfo is:%v", rollInfo)
	return rollInfo, nil
}
