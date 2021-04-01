// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"aliyun-ddns/config"
	"aliyun-ddns/domain"
	"aliyun-ddns/ip"
	"log"
	"os"
	"time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dns.aliyuncs.com")
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func createAliyunDnsClient(c *config.Configuration) (*alidns20150109.Client, error) {
	client, err := CreateClient(&c.DNSAuth.AppId, &c.DNSAuth.AppSecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	c := config.LoadConfig(os.Args)
	dnsmanager := domain.NewDomainManager(c)
	retryInterval := 60 * time.Second
	for {
		ipaddr, err := ip.GetIp()
		if err != nil {
			log.Print("get ip error", err)
			time.Sleep(retryInterval)
			retryInterval *= 2
			continue
		}
		dnsmanager.Update(ipaddr)
		time.Sleep(600 * time.Second)
		retryInterval = 60 * time.Second
	}
}
