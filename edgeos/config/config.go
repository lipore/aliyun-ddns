package config

import (
	"aliyun-ddns/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type EdgeosConfig struct {
	config.Configuration
	InterfaceMap map[string]string `yaml:"interface_map"`
}


func LoadConfig(configFile string) *EdgeosConfig {

	configReader, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	configBytes, err := ioutil.ReadAll(configReader)
	if err != nil {
		panic(err)
	}
	edgeConfig := EdgeosConfig{}
	yaml.Unmarshal(configBytes, &edgeConfig)
	config := config.Configuration{}
	yaml.Unmarshal(configBytes, &config)
	edgeConfig.Configuration = config
	return &edgeConfig
}