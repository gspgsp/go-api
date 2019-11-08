package cart

import (
	"edu_api/controllers"
	"edu_api/middlewares"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

type CartController struct {
	controller controllers.Controller
}

/**
添加购物车
*/
func (cart *CartController) AddCartInfo(w rest.ResponseWriter, r *rest.Request) {
	var addCart middlewares.AddCart
	if err := r.DecodeJsonPayload(&addCart); err != nil {
		log.Info("参数格式不正确:" + err.Error())
	}

	result, err := (&addCart).AddCartValidator()
	if err != nil {
		log.Info("验证错误:" + err.Error())
		cart.controller.Err = err
		cart.controller.JsonReturn(w, "result", err.Error())
		return
	}

	if result {
		code, message := cart.controller.BaseOrm.AddCartInfo(r, &addCart)
		if code == 0 {
			cart.controller.Err = nil
		} else {
			switch v := message.(type) {
			case string:
				cart.controller.Err = errors.New(v)
			}
		}

		cart.controller.JsonReturn(w, "cart_add", message)
	} else {
		cart.controller.Err = errors.New("未知错误")
		cart.controller.JsonReturn(w, "result", "")
	}
}
