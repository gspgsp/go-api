package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"strconv"
	"log"
)

func (baseOrm *BaseOrm) GetPlayList(r *rest.Request) (playLists []models.Chapter, err error) {

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		return
	}

	log.Printf("the id is:%v", id)

	return
}
