package cashier

import (
	"edu_api/controllers"
	"edu_api/middlewares"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

type CashierController struct {
	controller controllers.Controller
}

/**
生成支付信息
*/
func (cashier *CashierController) Payment(w rest.ResponseWriter, r *rest.Request) {
	var payment middlewares.Payment
	if err := r.DecodeJsonPayload(&payment); err != nil {
		log.Info("参数格式不正确:" + err.Error())
		return
	}

	result, err := (&payment).PaymentValidator()
	if result {
		code, message := cashier.controller.BaseOrm.Payment(r, &payment)
		if code == 0 {
			cashier.controller.Err = nil
		} else {
			switch v := message.(type) {
			case string:
				cashier.controller.Err = errors.New(v)
			case error:
				cashier.controller.Err = v
			}
		}
		cashier.controller.JsonReturn(w, "result", message)
	} else {
		cashier.controller.Err = err
		cashier.controller.JsonReturn(w, "result", "")
	}
}

func (cashier *CashierController) PayNotify(w rest.ResponseWriter, r *rest.Request) {
	var message string

	if r.PathParam("type") == "alipay" { //包括 余额、花呗(分期)
		message = cashier.controller.BaseOrm.PayNotify(r, "alipay")
	} else if r.PathParam("type") == "wechat_pay" { //包括h5支付、jsapi支付

	}

	//其实这里不能这样返回，因为支付宝那边只需要一个字符串success/fail，
	//后续准备把这个rest框架换为其它的，这个ResponseWriter没有
	//重写golang 自带的http.ResponseWriter的write()方法
	cashier.controller.JsonReturn(w, "result", message)

}
