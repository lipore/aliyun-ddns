// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"aliyun-ddns/config"
	"aliyun-ddns/domain/alidns"
	"aliyun-ddns/ip"
	"log"
	"os"
	"time"
)

func main() {
	c := config.LoadConfig(os.Args)
	dnsmanager := alidns.NewAlidnsDomainManager(c)
	for !c.DisableLoop {
		time.Sleep(600 * time.Second)
		ipaddr, err := ip.GetIp()
		if err != nil {
			log.Print("get ip error", err)
			continue
		}
		domains := dnsmanager.Domains()
		for _, domain := range *domains {
			domain.Update(ipaddr)
		}
	}
}
