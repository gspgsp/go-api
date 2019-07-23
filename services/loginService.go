package services

import (
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"edu_api/utils"
	"errors"
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

	return access_token, nil

}
