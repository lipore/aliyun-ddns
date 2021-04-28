package ip

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type IP string

func GetIp() (IP, error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("get ip error, status code: %d", res.StatusCode)
	}
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Printf("get ip: %s", ip)
	return IP(string(ip)), nil
}

func GetIp6() (IP, error) {
	res, err := http.Get("https://api64.ipify.org")
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("get ip error, status code: %d", res.StatusCode)
	}
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Printf("get ipv6: %s", ip)
	return IP(string(ip)), nil
}
