package order

import (
	"edu_api/controllers"
	"edu_api/middlewares"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

type OrderController struct {
	controller controllers.Controller
}

/**
提交订单
*/
func (order *OrderController) SubmitOrder(w rest.ResponseWriter, r *rest.Request) {
	var commitOrder middlewares.CommitOrder
	if err := r.DecodeJsonPayload(&commitOrder); err != nil {
		log.Info("参数格式不正确:" + err.Error())
	}

	result, err := (&commitOrder).CommitOrderValidator()

	if result {
		code, message := order.controller.BaseOrm.SubmitOrder(r, &commitOrder)
		if code == 0 {
			order.controller.Err = nil
		} else {
			switch v := message.(type) {
			case string:
				order.controller.Err = errors.New(v)
			}
		}
		order.controller.JsonReturn(w, "result", message)
	} else {
		order.controller.Err = err
		order.controller.JsonReturn(w, "result", "")
	}
}
