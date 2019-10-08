package services

import (
	"edu_api/models"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func (baseOrm *BaseOrm) GetExamRollTopicList(r *rest.Request) {
	var (
		rollList []models.RollModel
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
	test()

}

var (
	msgs           string      //信息
	msgChan        chan string //单条信息通道
	endTaskNumChan chan bool   //完成任务通知
	endTaskNum     int         //已完成任务数
)

func test() {
	msgs = "初始值\n"
	msgChan = make(chan string)
	endTaskNumChan = make(chan bool)
	endTaskNum = 0
	for i := 0; i < 5; i++ {
		go addMsg(i)
	}
L:
	//for  {
	//	fmt.Println("msgs\n", msgs)
	//	break L
	//}
	for {
		select {
		case msg, ok := <-msgChan: //获取单条信息
			if ok {
				msgs = fmt.Sprint(msgs, msg)
			}
		case <-endTaskNumChan: //获取处理完通知
			endTaskNum++
			if endTaskNum == 5 { //如果已完成任务等于总任务（退出这里是否有更好的办法？？）
				close(msgChan)
				break L
			}
		}
		fmt.Println("8888")
	}
	for i := 0; i < 10; i++ {
		fmt.Println("10086")
	}

	fmt.Printf("msg is:%v", msgs)
	fmt.Println("1")
	fmt.Println("2")
	fmt.Println("3")
}

func addMsg(i int) {
	defer func() {
		endTaskNumChan <- true
	}()
	//time.Sleep(1 * time.Second)
	msgChan <- fmt.Sprint("当前是新邮箱：", i)
}
