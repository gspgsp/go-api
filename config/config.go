package config

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var (
	root = flag.String("dir", "D:/gopath/src/edu_api", "设置配置文件的根路径")

	Config = parseYaml()
)

func parseYaml() *configuration {
	flag.Parse()

	cfg := new(configuration)
	cfg, err := cfg.yaml(*root + "/config/config.yaml")
	if err != nil {
		log.Info("parse yaml config error:", err.Error())
	}

	return cfg
}

type configuration struct {
	Queue  queue  `json:"queue",yaml:"queue"`
	Redis  redis  `json:"redis",yaml:"redis"`
	Alipay alipay `json:"alipay",yaml:"alipay"`
}

/**
对于只有一个元素的json tag标签可以解析
*/
type queue struct {
	Addr string `yaml:"addr"`
}

/**
有多个元素的时候，带有json tag的yaml标签无法解析；所以只能够去掉json tag
*/
type redis struct {
	CacheHost     string `yaml:"cache_host"`
	CachePort     int    `yaml:"cache_port"`
	CacheDatabase int    `yaml:"cache_database"`
	CacheUserName string `yaml:"cache_user_name"`
	CachePassword string `yaml:"cache_password"`
}

type alipay struct {
	AppId      string `yaml:"app_id"`
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
}

func (cfg *configuration) yaml(f string) (*configuration, error) {
	file, err := os.Open(f)
	if err != nil {
		log.Info("err is:", err.Error())
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Info("err is:", err.Error())
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Info("err is:", err.Error())
		return nil, err
	}

	return cfg, nil
}
