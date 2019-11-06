package home

import (
	"edu_api/controllers"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
)

type SlideController struct {
	controller controllers.Controller
}

/**
轮播信息
*/
func (slide *SlideController) GetSlide(w rest.ResponseWriter, r *rest.Request) {
	code, message := slide.controller.BaseOrm.GetSlide(r)
	if code == 0 {
		slide.controller.Err = nil
	} else {
		switch v := message.(type) {
		case string:
			slide.controller.Err = errors.New(v)
		}
	}

	slide.controller.JsonReturn(w, "slide", message)
}
