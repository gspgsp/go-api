package services

import (
	"edu_api/models"
	"github.com/iGoogle-ink/gopay"
)

var (
	user         models.User         //用户信息
	db           *BaseOrm            //数据库操作对象
	auth         string              //授权信
	aliPayClient *gopay.AliPayClient //支付宝支付客户端
	aliPayConfig AliPayConf          //支付宝支付配置
)
