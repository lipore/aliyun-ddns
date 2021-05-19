package ip

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type IP struct {
	V4, V6 string
}

func getIp(endpoint string) (string, error) {
	retryInterval := 30 * time.Second
	retryTimes := 0
	var res *http.Response
	for {
		if retryTimes > 5 {
			return "", fmt.Errorf("max retry time rearched")
		}
		var err error
		res, err = http.Get(endpoint)
		if err != nil || res.StatusCode != 200 {
			time.Sleep(retryInterval)
			retryTimes++
			continue
		}
		break
	}
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Printf("get ip: %s", ip)
	return string(ip), nil
}

func GetIp() (*IP, error) {
	var finalErr error = nil
	ipv4, err := getIp("https://api.ipify.org")
	if err != nil {
		ipv4 = ""
		finalErr = err
	}
	ipv6, err := getIp("https://api64.ipify.org")
	if err != nil {
		ipv6 = ""
		finalErr = err
	}
	return &IP{V4: ipv4, V6: ipv6}, finalErr
}
