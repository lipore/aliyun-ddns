package main

import (
	"aliyun-ddns/domain/alidns"
	"aliyun-ddns/edgeos/config"
	"aliyun-ddns/ip"
	"os"
)

func main() {
	cnf := config.LoadConfig("./config.yml")
	dnsDomainManager := alidns.NewAlidnsDomainManager(&cnf.Configuration)
	inter := os.Args[1]
	rr := cnf.InterfaceMap[inter]
	addr := ip.IP{V4: os.Args[2]}
	domains := dnsDomainManager.Domains()
	for _, domain := range *domains {
		if domain.RR() == rr {
			domain.Update(&addr)
		}
	}
}
