package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LogConfig struct {
	FileName string `yaml:"fileName"`
	MaxSize int `yaml:"maxSize"`
	MaxBackups int `yaml:"maxBackups"`
	MaxAge int `yaml:"maxAge"`
	Compress bool `yaml:"compress"`
	Level string `yaml:"level"`
}

type ConfigStruct struct {
	LogConfig `yaml:"log"`
	InternalConfig `yaml:"internal"`
	ClientConfig `yaml:"client"`
}

type InternalConfig struct {
	Port int `yaml:"port"`
	Ip string `yaml:"ip"`
}

type ClientConfig struct {
	Port int `yaml:"port"`
	Ip string `yaml:"ip"`
}

var Cfg = &ConfigStruct{}

func init(){
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic("读取配置文件失败！")
	}
	err = yaml.Unmarshal(data, Cfg)
	if err != nil {
		panic("转换配置失败！")
	}
}