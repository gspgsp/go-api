package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
)

/**
视频播放控制器
 */
type PlayController struct {
	controller controllers.Controller
}

func (play *PlayController) GetPlayList(w rest.ResponseWriter, r *rest.Request) {
	var playLists models.Media

	playLists, play.controller.Err = play.controller.BaseOrm.GetPlayList(r)
	play.controller.JsonReturn(w, play.controller, "playLists", playLists)
}
