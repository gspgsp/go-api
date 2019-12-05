package services

import (
	"edu_api/models"
	"edu_api/utils"
	"encoding/json"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strconv"
	"time"
)

/**
redis缓存类
*/
//redis配置对象
type RedisConf struct {
	Redis redisYaml
}

//redis配置
type redisYaml struct {
	CacheHost     string `yaml:"cache_host"`
	CachePort     string `yaml:"cache_port"`
	CacheDatabase string `yaml:"cache_database"`
	CacheUsername string `yaml:"cache_username"`
	CachePassword string `yaml:"cache_password"`
}

//获取redis配置对象
func getRedisConf() (redisConf RedisConf, err error) {
	conf := RedisConf{}
	cacheFile, err := ioutil.ReadFile("D:/gopath/src/edu_api/config/redis.yaml")

	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(cacheFile, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

//redis 线程池
func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		MaxActive:   1000,
		Dial: func() (redis.Conn, error) {
			conf, err := getRedisConf()
			if err != nil {
				//记录日志
				log.Printf("get conf error:%v", err)
			}

			c, err := redis.Dial("tcp", conf.Redis.CacheHost+":"+conf.Redis.CachePort)
			if err != nil {
				if c != nil {
					c.Close()
				}
				log.Printf("dial error:%v", err)
			}

			if _, err := c.Do("auth", conf.Redis.CachePassword); err != nil {
				if c != nil {
					c.Close()
				}
				log.Printf("auth error:%v", err)
			}

			if _, err := c.Do("select", conf.Redis.CacheDatabase); err != nil {
				if c != nil {
					c.Close()
				}
				log.Printf("select db error:%v", err)
			}

			return c, nil
		},
	}
}

//获取redis实例
func GetRedisConnection() redis.Conn {
	pool := redisPool()
	return pool.Get()
}

//获取redis缓存下的信息:
func GetRedisCache(authHeaderToken string, act, val string) (info string) {
	var j models.JwtClaim
	token := j.ParseToken(authHeaderToken)

	switch value := token.Claims.(jwt2.MapClaims)["id"].(type) {
	case float64:
		conn := GetRedisConnection()
		defer conn.Close()

		key := utils.ContactHashKey([]string{"user:", strconv.FormatFloat(value, 'f', -1, 64)}...)
		info, _ := redis.String(conn.Do(act, key, val))

		return info
	}

	return
}

/**
获取缓存的用户信息
*/
func GetUserInfo(authHeaderToken string) (user models.User) {
	info := GetRedisCache(authHeaderToken, "hget", "info")

	if err := json.Unmarshal([]byte(info), &user); err != nil {
		//记录日志
		log.Info("解析用户信息错误:", err.Error())
	}

	return
}

//设置用户全部课程最近观看课时记录(全局:包括最近观看课时，目前不用计算比例)
func SetLatestMediumPlayInfo(userId, lessonId interface{}) {
	conn := GetRedisConnection()
	defer conn.Close()

	conn.Do("hset", fmt.Sprintf(utils.LATEST_MEDIUM_PLAY, userId), utils.LATEST_LESION_INFO, lessonId)
}

func GetLatestMediumPlayInfo(userId interface{}) {

}

//设置某一课程最近观看课时记录(针对课程:包括最近观看课时、比例)
func SetClassLatestMediumPlayInfo(userId, courseId, info interface{}) {
	conn := GetRedisConnection()
	defer conn.Close()

	conn.Do("hset", fmt.Sprintf(utils.LATEST_MEDIUM_PLAY, userId), fmt.Sprintf(utils.LATEST_CLASS_WATCH_INFO, courseId), info)
}

//设置某一课程章节最近观看课时记录(针对章节:包括最近观看课时、比例)
func SetClassChapterMediumPlayInfo(userId, courseId, chapterId, info interface{}) {
	conn := GetRedisConnection()
	defer conn.Close()

	conn.Do("hset", fmt.Sprintf(utils.LATEST_MEDIUM_PLAY, userId), fmt.Sprintf(utils.LATEST_CLASS_CHAPTER_WATCH_INFO, courseId, chapterId), info)
}

/**
mongoDb缓存类
*/
