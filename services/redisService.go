package services

import (
	"flag"
	"github.com/garyburd/redigo/redis"
)

var (
	pool        *redis.Pool
	redisServer = flag.String("redisServer", ":6379", "")
)

type RedisPool struct {
}
