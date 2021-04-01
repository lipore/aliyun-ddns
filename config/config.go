package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	DNSAuth dnsAuth  `yaml:"dnsAuth"`
	Records []record `yaml:"records"`
}

type record struct {
	Domain string   `yaml:"domain"`
	RR     []string `yaml:"rr"`
}
type dnsAuth struct {
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

func LoadConfig(args []string) *Configuration {
	configFile := args[1]

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
