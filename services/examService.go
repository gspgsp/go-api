package services

import (
	"edu_api/models"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

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
	gradeChan <- grade
}
