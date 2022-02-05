package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	DNSAuth     dnsAuth  `yaml:"dnsAuth"`
	Records     []record `yaml:"records"`
	DisableLoop bool     `yaml:"disableLoop"`
}

type record struct {
	Domain string   `yaml:"domain"`
	RR     []string `yaml:"rr"`
	RRv6   []string `yaml:"rrv6"`
}
type dnsAuth struct {
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

func LoadConfig(configFile string) *Configuration {

	configReader, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	configBytes, err := ioutil.ReadAll(configReader)
	if err != nil {
		panic(err)
	}
	config := Configuration{}
	yaml.Unmarshal(configBytes, &config)
	return &config
}
