package main

import (
	"strconv"
)

type Service struct {
	ip       string
	port     int
	commands map[string]string
}

func (service Service) getIpAndPort() string {
	return service.ip + ":" + strconv.Itoa(service.port)
}
