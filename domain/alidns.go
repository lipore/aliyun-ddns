package domain

import (
	"aliyun-ddns/config"
	"aliyun-ddns/ip"
	"fmt"
	"log"
	"time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

func createClient(appid, appSecret string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &appid,
		// 您的AccessKey Secret
		AccessKeySecret: &appSecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dns.aliyuncs.com")
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

type AlidnsDomain struct {
	rr, domain, recordId string
	ip                   ip.IP
	domainManager        *AlidnsDomainManager
}

func (ad *AlidnsDomain) Update(ip ip.IP) error {
	if ip != ad.ip {
		return ad.updateDomainValue(ip)
	}
	return nil
}

func (ad *AlidnsDomain) updateDomainValue(ip ip.IP) error {
	retryInterval := 60 * time.Second
	for {
		client, err := createClient(ad.domainManager.appid, ad.domainManager.appSecret)
		if err != nil {
			return err
		}
		updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
			RR:       tea.String(ad.rr),
			RecordId: tea.String(ad.recordId),
			Type:     tea.String("A"),
			Value:    tea.String(string(ip)),
		}
		_, _err := client.UpdateDomainRecord(updateDomainRecordRequest)
		if _err != nil {
			fmt.Print(_err)
			log.Printf("update domain %s.%s to %s failed\n", ad.rr, ad.domain, ip)
			time.Sleep(retryInterval)
			log.Printf("retry\n")
			retryInterval *= 2
			continue
		}
		log.Printf("updated domain %s.%s to %s\n", ad.rr, ad.domain, ip)
		break
	}
	return nil
}

type AlidnsDomainManager struct {
	appid, appSecret string
	domains          []Domain
}

func newAlidnsDomainManager(c *config.Configuration) *AlidnsDomainManager {
	dm := AlidnsDomainManager{
		appid:     c.DNSAuth.AppId,
		appSecret: c.DNSAuth.AppSecret,
	}
	dm.domains = make([]Domain, 0)

	ipaddr, err := ip.GetIp()
	if err != nil {
		ipaddr = "127.0.0.2"
	}

	for _, record := range c.Records {
		for _, rr := range record.RR {
			dm.AddDomain(rr, record.Domain, ipaddr)
		}
	}
	return &dm
}

func (dm *AlidnsDomainManager) Update(ip ip.IP) error {
	for _, domain := range *dm.Domains() {
		err := domain.Update(ip)
		if err != nil {
			return err
		}
	}
	return nil
}
func (dm *AlidnsDomainManager) Domains() *[]Domain {
	return &dm.domains
}
func (dm *AlidnsDomainManager) AddDomain(rr, domain string, ipaddr ip.IP) {
	ad, err := dm.newAlidnsDomain(rr, domain, ipaddr)
	if err != nil {
		log.Print(err)
	}
	dm.domains = append(dm.domains, ad)
}

func (dm *AlidnsDomainManager) newAlidnsDomain(rr, domain string, ipaddr ip.IP) (ad *AlidnsDomain, err error) {
	client, err := createClient(dm.appid, dm.appSecret)
	if err != nil {
		return nil, err
	}
	describeSubDomainRecordsRequest := &alidns20150109.DescribeSubDomainRecordsRequest{
		SubDomain: tea.String(rr + "." + domain),
	}
	// 复制代码运行请自行打印 API 的返回值
	domainsResp, err := client.DescribeSubDomainRecords(describeSubDomainRecordsRequest)
	if err != nil {
		return nil, err
	}
	records := domainsResp.Body.DomainRecords.Record
	ad = &AlidnsDomain{
		rr:            rr,
		domain:        domain,
		domainManager: dm,
	}
	if len(records) != 0 {
		for _, record := range records {
			if *record.Type == "A" && *record.RR == rr {
				log.Printf("domain: %s.%s, found\n", rr, domain)
				ad.recordId = *record.RecordId
				ad.ip = ip.IP(*record.Value)
				return ad, err
			}
		}
	}
	log.Printf("domain: %s.%s, not found, will create it\n", rr, domain)

	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(domain),
		RR:         tea.String(rr),
		Type:       tea.String("A"),
		Value:      tea.String(string(ipaddr)),
		TTL:        tea.Int64(600),
	}
	addDomainResp, err := client.AddDomainRecord(addDomainRecordRequest)
	if err != nil {
		log.Printf("domain: %s.%s, create failed: %s \n", rr, domain, err.Error())
		return nil, err
	}
	ad.recordId = *addDomainResp.Body.RecordId
	ad.ip = ipaddr
	return ad, err

}
