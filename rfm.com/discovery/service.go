package main

import (
	"strconv"
)

type Service struct {
	ip      string
	port    int
	feature Feature
}

func (service Service) getIpAndPort() string {
	return service.ip + ":" + strconv.Itoa(service.port)
}
