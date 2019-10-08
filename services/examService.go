package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func (baseOrm *BaseOrm) GetExamRollTopicList(r *rest.Request) {
	var (
		rollList   []models.RollModel
		gradeChan  chan models.GradeModel
		endChan    chan bool
		endChanNum int
		allChanNum int
	)

	id, err := valid.ToInt(r.PathParam("id"))
	if err != nil {
		return
	}

	if err := baseOrm.GetDB().Table("h_exam_rolls").Where("course_id = ? and status = 2", id).Find(&rollList).Error; err != nil {
		log.Info("获取数据错误:" + err.Error())
		return
	}

	user = GetUserInfo(r.Header.Get("Authorization"))
	//计算是否答过当前考卷，以及考试结果
	if len(rollList) == 0 {
		return
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

				log.Printf("the grade is:%v", grade)
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
