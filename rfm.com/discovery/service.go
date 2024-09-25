package main

import (
	"net"
	"strconv"
)

type Service struct {
	ip      net.IP
	port    int
	feature Feature
}

func (service Service) getIpAndPort() string {
	return service.ip.String() + ":" + strconv.Itoa(service.port)
}
