package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/utils"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"time"
	"net/http"
	"bytes"
)

/**
VIP信息
*/
func (baseOrm *BaseOrm) GetVipInfo(r *rest.Request) (int, interface{}) {
	var (
		vip models.VipModel
	)
	id, err := valid.ToInt(r.PathParam("id"))
	if err != nil {
		log.Info("获取路由参数错误:" + err.Error())
		return 1, "获取路由参数错误:" + err.Error()
	}

	baseOrm.GetDB().Table("h_vips").Where("id = ?", id).First(&vip)
	user = GetUserInfo(r.Header.Get("Authorization"))

	if vip.ID == 0 {
		log.Info("会员信息不存在")
		return 1, "会员信息不存在"
	}

	if vip.ID > 0 && user.Id > 0 && user.Level == vip.VipLevel {
		vip.IsBuy = 1
	}

	vip.ParseCreatedAt, _ = FormatLocalTime(vip.CreatedAt)
	vip.ParseUpdatedAt, _ = FormatLocalTime(vip.UpdatedAt)

	return 0, vip
}

/**
创建会员订单
*/
func (baseOrm *BaseOrm) CreateVipOrder(r *rest.Request, vipOrder *middlewares.VipOrder) (int, string) {
	var vip models.VipModel
	if err := baseOrm.GetDB().Table("h_vips").Where("id = ?", vipOrder.Id).First(&vip).Error; err != nil {
		log.Info("获取数据错误:" + err.Error())
		return 1, err.Error()
	}

	user = GetUserInfo(r.Header.Get("Authorization"))
	if user.Level == vip.VipLevel {
		log.Info("您已经是荣耀终身会员，无需再次购买")
		return 1, "您已经是荣耀终身会员，无需再次购买"
	}
	var (
		vips      models.VipOrderModel
		createdAt string
		//updatedAt string
	)
	createdAt, _ = FormatLocalTime(time.Now())
	//updatedAt = createdAt
	baseOrm.GetDB().
		Table("h_vip_orders").
		Where("vip_id = ? and user_id = ? and status = 0 and payment_status = 0 and payment_expired_at > ?", vipOrder.Id, user.Id, createdAt).
		First(&vips)

	if vips.ID > 0 {
		log.Info("VIP订单已存在")
		return 1, "VIP订单已存在"
	}
	vipOrderData := make(map[string]interface{})

	vipOrderData["no"] = utils.GenerateOrderNo()
	vipOrderData["amount"] = vip.Price
	vipOrderData["source"] = vipOrder.Source
	vipOrderData["user_id"] = user.Id
	vipOrderData["vip_id"] = vip.ID
	vipOrderData["discount_amount"] = 0.0
	vipOrderData["created_at"] = createdAt
	vipOrderData["updated_at"] = createdAt
	t := time.Now().Unix()
	vipOrderData["payment_expired_at"] = fmt.Sprint(time.Unix(t+utils.PAYMENT_EXPIRED_HOUR*3600, 0).Format("2006-01-02 15:04:05"))

	if vip.DiscountEndAt > createdAt {
		vipOrderData["discount_amount"] = vip.Discount
		vipOrderData["payment_amount"] = fmt.Sprintf("%.2f", vip.Price-vip.Discount)
	} else {
		vipOrderData["payment_amount"] = vip.Price
	}

	insert_sql := "insert into `h_vip_orders` (`no`, `amount`, `source`, `user_id`, `vip_id`, `discount_amount`, `created_at`, `updated_at`, `payment_expired_at`, `payment_amount`) values"
	insert_value := fmt.Sprintf("(%s,%f,%s,%d,%d,%f,%s,%s,%s,%f)", vipOrderData["no"], vipOrderData["amount"], `'`+vipOrder.Source+`'`, vipOrderData["user_id"], vipOrderData["vip_id"], vipOrderData["discount_amount"], "'"+createdAt+"'", "'"+createdAt+"'", "'"+fmt.Sprint(time.Unix(t+utils.PAYMENT_EXPIRED_HOUR*3600, 0).Format("2006-01-02 15:04:05"))+"'", vipOrderData["payment_amount"])

	tx := baseOrm.GetDB().Begin()
	err := tx.Exec(insert_sql + insert_value).Error
	if err != nil {
		log.Info("事务操作出错:" + fmt.Sprintf("插入VIP订单错误:%s", err.Error))
		tx.Rollback()
		return 1, "VIP订单创建失败"
	} else {
		log.Info("插入VIP订单成功")

		//向任务队列插入任务
		//time.AfterFunc(time.Second*3600*48, func() {
		//	log.Info("任务完成，订单号为:"+fmt.Sprintf("%s", vipOrderData["no"]))
		//})
		SendDelayQueueRequest()
		tx.Commit()
		return 0, "VIP订单创建成功"
	}
}

/**
发送任务
 */
func SendDelayQueueRequest()  {
	post :=`{"topic":"order","id":"15702398324","delay":3600,"ttr":120,"body":"do something"}`
	var jsonStr = []byte(post)

	url := "http://39.106.141.78:9277/push"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		//panic(err)
		fmt.Printf("the err is:%v", err.Error())
	}
	defer resp.Body.Close()
}
