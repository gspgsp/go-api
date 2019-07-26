package middlewares

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
)

/**
token中间件
 */
type AuthTokenMiddleware struct {

}

//实现一下HandlerFunc，和http.HandlerFunc是一样的效果，其实这个中间件的就像执行了一个匿名函数一样
func (atm *AuthTokenMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	//这里可以执行atm的其它操作
	log.Println("the atm middle")

	return func(writer rest.ResponseWriter, request *rest.Request) {
		//前置处理
		authHeader := request.Header.Get("Authorization")
		log.Printf("the auth header is:%v", authHeader)

		//相当于next()
		handler(writer, request)
	}
}
