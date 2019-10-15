package vip

import (
	"edu_api/controllers"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
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

}
