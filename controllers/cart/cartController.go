package cart

import (
	"edu_api/controllers"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
)

type CartController struct {
	controller controllers.Controller
}

func (cart *CartController) AddCartInfo(w rest.ResponseWriter, r *rest.Request) {
	code, message := cart.controller.BaseOrm.AddCartInfo(r)
	if code == 0 {
		cart.controller.Err = nil
	} else {
		switch v := message.(type) {
		case string:
			cart.controller.Err = errors.New(v)
		}
	}

	cart.controller.JsonReturn(w, "cart_add", message)
}
