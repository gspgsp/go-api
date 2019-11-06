package middlewares

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"regexp"
)

/**
路由判断中间件
*/
func IfMiddleware() *rest.IfMiddleware {

	middleware := &rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {

			path := request.URL.Path

			expr := `(/login)|(/register)|(/package)|(/course[/\d+]?)|(/category)|(/chapter[/\d+]?)|(/lecture[/\d+]?)|(/review[/\d+]?)|(/recommend[/\d+]?)|(/try[/\d+]?)|(/compose[/\d+]?)|(/package[/\d+]?)|(/notice)|(/slide)`
			re, _ := regexp.Compile(expr)

			all := re.FindAllString(path, -1)

			for _, item := range all {
				log.Printf("the item is:%v", string(item))
				if len(string(item)) > 0 {
					return false
				}
			}

			return true
		},
		IfTrue: new(AuthTokenMiddleware), //或者&middlewares.AuthTokenMiddleware{}
	}

	return middleware
}
