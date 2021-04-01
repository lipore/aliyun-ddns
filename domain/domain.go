package domain

import (
	"aliyun-ddns/config"
	"aliyun-ddns/ip"
)

type Domain interface {
	Update(ip.IP) error
}

type DomainManager interface {
	Update(ip.IP) error
	Domains() *[]Domain
	AddDomain(rr, domain string, ipaddr ip.IP)
}

func NewDomainManager(c *config.Configuration) DomainManager {
	domainManager := newAlidnsDomainManager(c)
	return domainManager
}
