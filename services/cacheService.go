package services

import (
	"flag"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v3"
)

/**
redis cache
 */
var (
	pool        *redis.Pool
	redisServer = flag.String("redisServer", ":6379", "")
)

/**
redis 线程池
 */
type RedisPool struct {
}

/**
redis配置对象
 */
type RedisConf struct {
	Redis redisYaml
}

/**
redis配置
 */
type redisYaml struct {
	CacheHost     string `yaml:"cache_host"`
	CachePort     string `yaml:"cache_port"`
	CacheDatabase string `yaml:"cache_database"`
	CacheUsername string `yaml:"cache_username"`
	CachePassword string `yaml:"cache_password"`
}

/**
获取redis配置对象
 */
func GetRedisConf() {
	conf := RedisConf{}
	cacheFile, err := ioutil.ReadFile("E:/GoProjects/src/edu_api/config/redis.yaml")

	if err != nil {
		log.Printf("the yaml error is:%v", err)
	}

	err = yaml.Unmarshal(cacheFile, &conf)
	if err != nil {
		log.Printf("the yaml unmarshal error is:%v", err)
	}

	log.Print("the cache conf is:\n", conf.Redis)
}
