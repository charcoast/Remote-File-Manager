package main

import (
	"rfm.com/commom"
	"strconv"
)

type Service struct {
	ip      string
	port    int
	feature commom.Feature
}

func (service Service) getIpAndPort() string {
	return service.ip + ":" + strconv.Itoa(service.port)
}
