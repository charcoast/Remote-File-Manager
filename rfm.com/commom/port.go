package commom

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func IsPortFree(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	err = listener.Close()
	if err != nil {
		return false
	}
	return true
}

func RandomPort() int {
	// Define o intervalo de portas
	minPort := 1024
	maxPort := 65535

	// Sementeia o gerador de números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Gera um número de porta aleatório
	return rand.Intn(maxPort-minPort+1) + minPort
}

func GetRandomPort() int {
	for {
		port := RandomPort()
		if IsPortFree(port) {
			return port
		}
	}
}
