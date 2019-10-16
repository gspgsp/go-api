package vip

import (
	"edu_api/controllers"
	"edu_api/middlewares"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

type VipController struct {
	controller controllers.Controller
}

/**
VIP信息
*/
func (vip *VipController) GetVipInfo(w rest.ResponseWriter, r *rest.Request) {
	code, message := vip.controller.BaseOrm.GetVipInfo(r)
	if code == 0 {
		vip.controller.Err = nil
	} else {
		switch v := message.(type) {
		case string:
			vip.controller.Err = errors.New(v)
		}
	}

	vip.controller.JsonReturn(w, "vip_info", message)
}

/**
创建会员订单
*/
func (vip *VipController) CreateVipOrder(w rest.ResponseWriter, r *rest.Request) {
	var vipOrder middlewares.VipOrder
	if err := r.DecodeJsonPayload(&vipOrder); err != nil {
		log.Info("参数格式不正确:" + err.Error())
	}

	result, err := (&vipOrder).VipOrderValidator()
	if err != nil {
		log.Info("验证错误:" + err.Error())
		vip.controller.Err = err
		vip.controller.JsonReturn(w, "result", err.Error())
		return
	}

	if result {
		vip.controller.BaseOrm.CreateVipOrder(r, &vipOrder)
	} else {
		vip.controller.Err = errors.New("未知错误")
		vip.controller.JsonReturn(w, "result", "")
	}
}
