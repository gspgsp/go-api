package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"strconv"
)

/**
轮播信息
*/
func (baseOrm *BaseOrm) GetSlide(r *rest.Request) (int, interface{}) {
	params := r.URL.Query()
	port, _ := strconv.Atoi(params.Get("port"))
	limit, _ := strconv.Atoi(params.Get("limit"))
	where := map[string]interface{}{"status": 1}

	if port > 0 && valid.InRange(port, 1, 2) == false {
		return 1, "参数错误"
	} else if port > 0 {
		where["port"] = port
	}

	if limit == 0 {
		limit = 10
	}

	var slides []models.SlideModel

	if err := baseOrm.GetDB().Table("h_slides").Select("id, port, title, url, carousel, description").Where(where).Limit(limit).Find(&slides).Error; err != nil {
		log.Info("获取公告出错")
		return 1, err.Error()
	}

	return 0, slides
}
