package middlewares

import (
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"net/http"
	"edu_api/models"
)

/**
token中间件
 */
type AuthTokenMiddleware struct {
}

/**
实现一下HandlerFunc，和http.HandlerFunc是一样的效果，其实这个中间件的就像执行了一个匿名函数一样
MiddlewareFunc会像调用构造函数一样被调用，AuthTokenMiddleware初始化的时候会调用这个方法
 */
func (atm *AuthTokenMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	//这里可以执行atm的其它操作
	log.Info("the atm middle")

	return func(writer rest.ResponseWriter, request *rest.Request) {
		//前置处理
		authHeaderToken := request.Header.Get("Authorization")

		//未设置token
		if authHeaderToken == "" {
			atm.unauthorized(writer, "请登录授权之后，再请求")
			return
		}

		//token不正确
		var j models.JwtClaim
		if _, err := j.VerifyToken(authHeaderToken); err != nil {
			atm.unauthorized(writer, "授权信息不正确，请重新授权")
			return
		}

		//相当于next()
		handler(writer, request)
	}
}

func (atm *AuthTokenMiddleware) unauthorized(writer rest.ResponseWriter, msg string) {
	writer.Header().Set("www-authenticate", "Token realm= jwt auth")
	rest.Error(writer, msg, http.StatusUnauthorized)
}
