package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var DefaultAPIConf = APIConfig{
	DSN: "root:root123@tcp(localhost:3306)/nomo?charset=utf8mb4&parseTime=True&loc=Local",
}

type APIConfig struct {
	Filename string `yaml:"-" json:"-"`
	DSN      string `json:"dsn" yaml:"dsn"`
}

func LoadAPIConfig(f string) (APIConfig, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return DefaultAPIConf, err
	}

	var c APIConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return DefaultAPIConf, err
	}

	c.Filename = f
	return c, nil
}
