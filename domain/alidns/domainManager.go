package alidns

import (
	"aliyun-ddns/config"
	"aliyun-ddns/domain"
	"aliyun-ddns/ip"
	"log"
)

type AlidnsDomainManager struct {
	appid, appSecret string
	domains          []domain.Domain
}

func NewAlidnsDomainManager(c *config.Configuration) *AlidnsDomainManager {
	dm := AlidnsDomainManager{
		appid:     c.DNSAuth.AppId,
		appSecret: c.DNSAuth.AppSecret,
	}
	dm.domains = make([]domain.Domain, 0)

	ipaddr := &ip.IP{
		V4: "127.0.0.2",
		V6: "::1",
	}
	for _, record := range c.Records {
		domainType := "A"
		for _, rr := range record.RR {
			dm.AddDomain(rr, record.Domain, domainType, ipaddr)
		}
		domainType = "AAAA"
		for _, rr := range record.RRv6 {
			dm.AddDomain(rr, record.Domain, domainType, ipaddr)
		}
	}
	return &dm
}

func (dm *AlidnsDomainManager) Domains() *[]domain.Domain {
	return &dm.domains
}
func (dm *AlidnsDomainManager) AddDomain(rr, domain, domainType string, ipaddr *ip.IP) {
	ad, err := dm.findOrCreateDomain(rr, domain, ipaddr, domainType)
	if err != nil {
		log.Print(err)
		return
	}
	dm.domains = append(dm.domains, ad)
}

func (dm *AlidnsDomainManager) findOrCreateDomain(rr, domain string, ipaddr *ip.IP, domainType string) (ad *AlidnsDomain, err error) {
	ipvalue := ipByDomainType(ipaddr, domainType)
	client, err := createClient(dm.appid, dm.appSecret)
	if err != nil {
		return nil, err
	}
	ad, err = findDomain(rr, domain, domainType, client, dm)
	if err != nil {
		return nil, err
	}

	if ad != nil {
		return ad, nil
	}
	log.Printf("domain: %s.%s, not found, will create it\n", rr, domain)
	ad, err = createDomain(rr, domain, domainType, ipvalue, client, dm)
	return ad, err

}
