package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"edu_api/utils"
	"strings"
	"errors"
)

type JwtClaim struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

var (
	Secret     = utils.TOKEN_SECRET
	ExpireTime = 3600
)

/**
设置过期时间
 */
func (j *JwtClaim) SetExpireAt(expireAt int64) {
	j.ExpiresAt = expireAt
}

/**
设置使用时间
 */
func (j *JwtClaim) SetIssUsedAt(issused int64) {
	j.IssuedAt = issused
}

/**
生成token
 */
func (j JwtClaim) AccessToken() (accessToken string, err error) {
	j.SetExpireAt(time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix())
	j.SetIssUsedAt(time.Now().Unix() - 1)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	accessToken, err = token.SignedString([]byte(Secret)) //这里虽然是interface类型，但是实际上要传[]byte

	return
}

/**
验证token
 */
func (j JwtClaim) VerifyToken(accessToken string) (token *jwt.Token, err error) {

	tokenArr := strings.Split(accessToken, " ")

	if length := len(tokenArr); length != 2 {
		return nil, errors.New(utils.TOKEN_PARAM_REQUIRE)
	}

	if tokenArr[0] != utils.TOKEN_TYPE {
		return nil, errors.New(utils.TOKEN_TYPE_ERROR)
	}

	//不要用ParseWithClaims，因为不好取返回值
	token, err = jwt.Parse(tokenArr[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})

	if err != nil {
		return nil, errors.New(utils.TOKEN_PARSE_ERROR)
	}

	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(utils.TOKEN_INVDLID)
	}

	return token, nil
}

/**
直接解析accessToken
 */
func (j JwtClaim) ParseToken(accessToken string) (token *jwt.Token) {
	tokenArr := strings.Split(accessToken, " ")
	token, _ = jwt.Parse(tokenArr[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})

	return token
}
