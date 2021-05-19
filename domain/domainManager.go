package domain

import (
	"aliyun-ddns/ip"
)

type DomainManager interface {
	Domains() *[]Domain
	AddDomain(rr, domain, domainType string, ipaddr ip.IP)
}
