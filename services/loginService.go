package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"edu_api/utils"
	"errors"
	"strconv"
	"log"
	"encoding/json"
)

var (
	login Login
	user  models.User
	jwt   models.JwtClaim
)

type Login struct {
	Mobile   string
	Password string
}

/**
登录
 */
func (baseOrm *BaseOrm) Login(r *rest.Request) (access_token string, err error) {

	if err = r.DecodeJsonPayload(&login); err != nil {
		return "", err
	}

	if err = baseOrm.GetDB().Table("h_users").Where("mobile = ?", login.Mobile).Find(&user).Error; err != nil {
		return "", err
	}

	if res := utils.PasswordVerify(login.Password, user.Password); res == false {
		return "", errors.New(utils.LOGIN_ERROR_NOTICE)
	}

	//准备access_token
	jwt.Id = user.Id
	if access_token, err = jwt.AccessToken(); err != nil {
		return "", err
	}

	//将用户信息缓存到redis
	conn := GetRedisConnection()
	defer conn.Close()
	key := utils.ContactHashKey([]string{"user:", strconv.Itoa(user.Id)}...)

	//用户信息(注册信息)
	jsonUserInfo, _ := json.Marshal(user)
	if v, err := conn.Do("hsetnx", key, "info", jsonUserInfo); v == 0 {
		//记录日志
		log.Printf("the hsetnx err is:%v", err)
	}

	//用户详细信息

	//用户认证信息

	return access_token, nil

}
