package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func (baseOrm *BaseOrm) GetVipInfo(r *rest.Request) {
	id, err := valid.ToInt(r.PathParam("id"))
	if err != nil {
		log.Info("获取路由参数错误")
	}
}
