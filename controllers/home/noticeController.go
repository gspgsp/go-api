package home

import (
	"edu_api/controllers"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
)

type NoticeController struct {
	controller controllers.Controller
}

/**
公告信息
*/
func (notice *NoticeController) GetNotice(w rest.ResponseWriter, r *rest.Request) {
	code, message := notice.controller.BaseOrm.GetNotice(r)
	if code == 0 {
		notice.controller.Err = nil
	} else {
		switch v := message.(type) {
		case string:
			notice.controller.Err = errors.New(v)
		}
	}

	notice.controller.JsonReturn(w, "notice", message)
}
