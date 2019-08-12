package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"strconv"
	"log"
)

func (baseOrm *BaseOrm) GetPlayList(r *rest.Request) (playLists []models.Media, err error) {

	id, err := strconv.Atoi(r.PathParam("id"))
	lesionId, err := strconv.Atoi(r.PathParam("lesion_id"))
	if err != nil {
		return
	}

	chapterList := make([]models.Chapter, 10)

	first := make([]models.Chapter, 1)

	if err := baseOrm.GetDB().Table("h_edu_chapters").Where("course_id = ? and status = 'published'", id).Find(&chapterList).Error; err != nil {
		return nil, err
	}

	if err := baseOrm.GetDB().Table("h_edu_chapters").Where("id = ? and status = 'published'", lesionId).Find(&first).Error; err != nil {
		return nil, err
	}

	if len(first) == 0 {
		//就选择当前课程下最远的一条数据
		if err := baseOrm.GetDB().Table("h_edu_chapters").Where("course_id = ? and status = 'published' and type = 'lesson'", id).Order("updated_at asc").Order("created_at asc").Find(&first).Error; err != nil {
			return nil, err
		}

	}

	log.Printf("the playLists is:%v", first)

	return
}
