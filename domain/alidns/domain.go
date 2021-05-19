package alidns

import (
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
	ip                   string
	domainManager        *AlidnsDomainManager
	domainType           string
}

func (ad *AlidnsDomain) Update(ip *ip.IP) error {
	ipvalue := ipByDomainType(ip, ad.domainType)
	client, err := createClient(ad.domainManager.appid, ad.domainManager.appSecret)
	if err != nil {
		return err
	}
	domain, err := findDomain(ad.rr, ad.domain, ad.domainType, client, ad.domainManager)
	if err != nil {
		// online check failed, fallback to offline
		domain = ad
	}
	if ipvalue != domain.ip {
		return ad.updateDomainValue(ipvalue, client)
	}
	return nil
}

func (ad *AlidnsDomain) updateDomainValue(ipvalue string, client *alidns20150109.Client) error {
	retryInterval := 60 * time.Second
	for {
		updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
			RR:       tea.String(ad.rr),
			RecordId: tea.String(ad.recordId),
			Type:     tea.String(ad.domainType),
			Value:    tea.String(ipvalue),
		}
		_, _err := client.UpdateDomainRecord(updateDomainRecordRequest)
		if _err != nil {
			fmt.Print(_err)
			log.Printf("update domain %s.%s to %s failed\n", ad.rr, ad.domain, ipvalue)
			time.Sleep(retryInterval)
			log.Printf("retry\n")
			retryInterval *= 2
			continue
		}
		log.Printf("updated domain %s.%s to %s\n", ad.rr, ad.domain, ipvalue)
		break
	}
	return nil
}

func ipByDomainType(ip *ip.IP, domainType string) string {
	ipvalue := ""
	if domainType == "A" {
		ipvalue = ip.V4
	} else if domainType == "AAAA" {
		ipvalue = ip.V6
	}
	return ipvalue
}

func findDomain(rr, domain, domainType string, client *alidns20150109.Client, dm *AlidnsDomainManager) (*AlidnsDomain, error) {
	describeSubDomainRecordsRequest := &alidns20150109.DescribeSubDomainRecordsRequest{
		SubDomain: tea.String(rr + "." + domain),
	}
	// 复制代码运行请自行打印 API 的返回值
	domainsResp, err := client.DescribeSubDomainRecords(describeSubDomainRecordsRequest)
	if err != nil {
		return nil, err
	}
	records := domainsResp.Body.DomainRecords.Record
	if len(records) != 0 {
		ad := &AlidnsDomain{
			rr:            rr,
			domain:        domain,
			domainManager: dm,
		}
		for _, record := range records {
			if *record.RR == rr && *record.Type == domainType {
				log.Printf("domain: %s.%s, found\n", rr, domain)
				ad.recordId = *record.RecordId
				ad.ip = *record.Value
				ad.domainType = *record.Type
				return ad, nil
			}
		}
	}
	return nil, nil
}

func createDomain(rr, domain, domainType, ipvalue string, client *alidns20150109.Client, dm *AlidnsDomainManager) (*AlidnsDomain, error) {

	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(domain),
		RR:         tea.String(rr),
		Type:       tea.String(domainType),
		Value:      tea.String(ipvalue),
		TTL:        tea.Int64(600),
	}
	addDomainResp, err := client.AddDomainRecord(addDomainRecordRequest)
	if err != nil {
		log.Printf("domain: %s.%s, create failed: %s \n", rr, domain, err.Error())
		return nil, err
	}

	return &AlidnsDomain{
		rr:            rr,
		domain:        domain,
		domainManager: dm,
		recordId:      *addDomainResp.Body.RecordId,
		ip:            ipvalue,
		domainType:    domainType,
	}, nil
}
