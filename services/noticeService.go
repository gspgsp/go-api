package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

/**
公告信息
*/
func (baseOrm *BaseOrm) GetNotice(r *rest.Request) (int, interface{}) {
	params := r.URL.Query()
	nType, _ := strconv.Atoi(params.Get("type"))
	limit, _ := strconv.Atoi(params.Get("limit"))

	if valid.InRange(nType, 1, 3) == false {
		return 1, "参数错误"
	}

	if limit == 0 {
		limit = 10
	}

	var notices []models.NoticeModel
	if err := baseOrm.GetDB().Table("h_notices").Where("status = 1 and type = ? and start_at < ? and (end_at > ? or end_at = 0)", nType, int32(time.Now().Unix()), int32(time.Now().Unix())).Select("id, type, title, content, status").Limit(limit).Find(&notices).Error; err != nil {
		log.Info("获取公告出错")
		return 1, err.Error()
	}

	return 0, notices
}
