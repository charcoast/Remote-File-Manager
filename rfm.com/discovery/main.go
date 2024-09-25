package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const DiscoveryPort = 7070
const ClientPort = 9090

var services []Service

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go listenToClient(&wg)
	go listenToServices(&wg)

	fmt.Println("****RUNNING****")

	wg.Wait()
}

func listenToServices(wg *sync.WaitGroup) error {
	serviceListener, err := net.Listen("tcp", ":"+strconv.Itoa(DiscoveryPort))
	defer serviceListener.Close()
	defer wg.Done()
	if err != nil {
		return err
	}

	for {
		conn, err := serviceListener.Accept()
		if err != nil {
			fmt.Println("ERRO AO INICIAR CONEXÃO COM UM SERVIÇO")
			fmt.Println(err)
			continue
		}

		go handleServiceDiscovery(conn)
	}

}

func handleServiceDiscovery(conn net.Conn) error {
	defer conn.Close()

	netData, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		return err
	}

	featureRegister := FeatureRegister{}

	err = json.Unmarshal([]byte(netData), &featureRegister)
	if err != nil {
		return err
	}

	var addr = conn.RemoteAddr().(*net.TCPAddr)

	service := Service{ip: addr.IP, port: featureRegister.Port, feature: Feature{Prefixes: featureRegister.Prefixes}}

	services = append(services, service)
	fmt.Printf("\n\nServiço adicionado. IP/Porta: %s Prefixos: %s", service.getIpAndPort(), strings.Join(service.feature.Prefixes, ","))

	return nil
}

func listenToClient(wg *sync.WaitGroup) error {
	clientListener, err := net.Listen("tcp", ":"+strconv.Itoa(ClientPort))
	defer clientListener.Close()
	defer wg.Done()

	if err != nil {
		return err
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

}

func handleClientCommand(conn net.Conn) error {

	defer conn.Close()
	netData, err := bufio.NewReader(conn).ReadString('\n')

	fmt.Println("RECEBEU O COMANDO: " + netData)

	if err != nil {
		return err
	}

	var index = slices.IndexFunc(services, func(s Service) bool {
		return slices.Contains(s.feature.Prefixes, netData)
	})

	if index == -1 {
		fmt.Fprintf(conn, "NÃO HÁ SERVIÇO CAPAZ DE RESPONDER SUA SOLICITAÇÃO")
		return nil
	}

	var service Service = services[index]

	c, err := net.Dial("tcp", service.getIpAndPort())
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Fprintf(c, netData) // Envia o texto pela conexão

	message, _ := bufio.NewReader(c).ReadString('\n') // Aguarda resposta do servidor
	fmt.Print("RESPONDEU: " + message)

	return nil
}
