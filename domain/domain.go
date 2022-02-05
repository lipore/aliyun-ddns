package domain

import (
	"aliyun-ddns/ip"
)

type Domain interface {
	Update(*ip.IP) error
	RR() string
}
