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

	log.Info("cfg is:", cfg.Queue.Addr)
	return cfg
}

type configuration struct {
	Queue queue `json:"queue",yaml:"queue"`
}

type queue struct {
	Addr string `json:"addr",yaml:"addr"`
}

func (cfg *configuration) yaml(f string) (*configuration, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
