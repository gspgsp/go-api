package order

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
)

type OrderController struct {
	controller controllers.Controller
}

func (order *OrderController) SubmitOrder(w rest.ResponseWriter, r *rest.Request) {

}
