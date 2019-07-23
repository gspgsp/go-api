package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtClaim struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

var (
	Secret     = "1qaz2wsx"
	ExpireTime = 3600
)

/**
设置过期时间
 */
func (j *JwtClaim) SetExpireAt(expireAt int64) {
	j.ExpiresAt = expireAt
}

func (j *JwtClaim) SetIssUsedAt(issused int64)  {
	j.IssuedAt = issused
}

func (j JwtClaim) AccessToken() (accessToken string, err error)  {
	j.SetExpireAt(time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix())
	j.SetIssUsedAt(time.Now().Unix() - 1)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	accessToken, err = token.SignedString([]byte(Secret))//这里虽然是interface类型，但是实际上要传[]byte

	return
}
