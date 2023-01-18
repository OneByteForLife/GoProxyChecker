package models

import (
	"net"
	"time"
)

type ProxyList struct {
	ID        int
	Type      string
	IP        string
	Port      int
	Speed     int
	AnonLVL   string
	City      string
	Country   string
	LastCheck time.Time
}

type NetClient interface {
	Check(val string, proxy func(string, string) (net.Conn, error)) (ProxyRespons, error)
}

type ProxyRespons interface {
	IsSuccess() bool
	GetStatusCodeRaw() int
}
