package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rfm.com/commom"
	"slices"
	"strconv"
)

const DiscoveryPort = 7070
const ClientPort = 9090

var services []Service

func main() {

	listenToClient()
	listenToServices()

	fmt.Println("****RUNNING****")

}

func listenToServices() {

	router := http.NewServeMux()

	router.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var featureRegister commom.FeatureRegister
		_ = json.NewDecoder(r.Body).Decode(&featureRegister)
		statusCode, message := handleServiceDiscovery(r.RemoteAddr, featureRegister)

		w.WriteHeader(statusCode)
		w.Write([]byte(message))
	})

	http.ListenAndServe(":"+strconv.Itoa(DiscoveryPort), router)

}

func handleServiceDiscovery(addr string, featureRegister commom.FeatureRegister) (int, string) {

	service := Service{ip: addr, port: featureRegister.Port, commands: featureRegister.Commands}

	services = append(services, service)
	message := fmt.Sprintf("\n\nServiço adicionado. IP/Porta: %s Prefixos: %s", service.getIpAndPort(), service.commands)
	fmt.Println(message)
	return 200, message
}

func listenToClient() {
	router := http.NewServeMux()

	router.HandleFunc("POST /command", func(w http.ResponseWriter, r *http.Request) {
		var command commom.Command
		_ = json.NewDecoder(r.Body).Decode(&command)
		statusCode, message := handleClientCommand(command.Command)

		w.WriteHeader(statusCode)
		w.Write([]byte(message))
	})

	http.ListenAndServe(":"+strconv.Itoa(ClientPort), router)
}

func handleClientCommand(command string) (int, string) {
	fmt.Println("RECEBEU O COMANDO: " + command)

	var index = slices.IndexFunc(services, func(s Service) bool {
		_, ok := s.commands[command]
		return ok
	})

	if index == -1 {
		return 500, "NÃO HÁ SERVIÇO CAPAZ DE RESPONDER SUA SOLICITAÇÃO"
	}

	service := services[index]
	url := fmt.Sprintf("http://%s:%s/", service.ip, service.port)
	body, err := json.Marshal(commom.Command{Command: command})

	if err != nil {
		return 400, "FALHA AO DECODIFICAR COMANDO"
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return 500, "NÃO FOI POSSÍVEL EXECUTAR O COMANDO"
	}

	var result string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	fmt.Print("RESPONDEU: " + result)

	return 200, result
}
