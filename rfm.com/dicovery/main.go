package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

const DISCOVERY_PORT = 7070
const CLIENT_PORT = 9090

var services []Service

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go listenToClient(&wg)
	go listenToServices(&wg)

	wg.Wait()
}

func listenToServices(wg *sync.WaitGroup) {
	serviceListener, err := net.Listen("tcp", ":"+strconv.Itoa(DISCOVERY_PORT))
	defer serviceListener.Close()
	if err != nil {
		panic(err)
	}

	for {
		conn, err := serviceListener.Accept()
		if err != nil {
			fmt.Println("ERRO AO INICIAR CONEXÃO COM UM SERVIÇO")
			fmt.Println(err)
			break
		}

		go handleServiceDiscovery(conn)
	}

	wg.Done()
}

func handleServiceDiscovery(conn net.Conn) {
	defer conn.Close()

	netData, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	featureRegister := FeatureRegister{}

	err = json.Unmarshal([]byte(netData), &featureRegister)
	if err != nil {
		return
	}

	var addr = conn.RemoteAddr().(*net.TCPAddr)

	service := Service{ip: addr.IP, port: featureRegister.Port, feature: Feature{Prefixes: featureRegister.Prefixes}}

	services = append(services, service)
	fmt.Printf("\n\nServiço adicionado. IP/Porta: %s Prefixos: %s", service.getIpAndPort(), strings.Join(service.feature.Prefixes, ","))
}

func listenToClient(wg *sync.WaitGroup) {
	clientListener, err := net.Listen("tcp", ":"+strconv.Itoa(CLIENT_PORT))
	defer clientListener.Close()

	if err != nil {
		panic(err)
	}

	for {
		conn, err := clientListener.Accept()
		if err != nil {
			fmt.Println("ERRO AO INICIAR CONEXÃO COM CLIENTE")
			fmt.Println(err)
			continue
		}

		go handleClientCommand(conn)
	}

	wg.Done()
}

func handleClientCommand(conn net.Conn) {
	fmt.Println("Executando conexão do cliente")
}
