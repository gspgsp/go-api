package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
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
func (baseOrm *BaseOrm) CreateVipOrder(r *rest.Request) {

}
