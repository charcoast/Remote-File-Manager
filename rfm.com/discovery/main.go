package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"github.com/thoas/go-funk"
	"io"
	"net"
	"net/http"
	"rfm.com/common"
	"slices"
	"strconv"
	"strings"
)

const DiscoveryPort = 7070
const ClientPort = 9090

var services []Service
var c *cors.Cors

func main() {

	fmt.Println("****RUNNING****")
	c = cors.New(cors.Options{AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
		AllowCredentials: true,
		Debug:            true})

	go listenToServices()
	listenToClient()
}

func getCors() *cors.Cors {
	return cors.New(cors.Options{AllowedOrigins: []string{"localhost:5173"},
		AllowCredentials: true,
		Debug:            true})
}

func listenToServices() {
	router := http.NewServeMux()
	router.HandleFunc("POST /register", handleRegisterService)
	http.ListenAndServe(":"+strconv.Itoa(DiscoveryPort), c.Handler(router))
}

func handleRegisterService(w http.ResponseWriter, r *http.Request) {
	var featureRegister common.FeatureRegister
	_ = json.NewDecoder(r.Body).Decode(&featureRegister)
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	statusCode, message := handleServiceDiscovery(host, featureRegister)

	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func handleServiceDiscovery(addr string, featureRegister common.FeatureRegister) (int, string) {

	featureRegister.Commands = funk.Map(featureRegister.Commands, func(k string, v string) (string, string) {

		return k, strings.TrimPrefix(v, "/")
	}).(map[string]string)

	service := Service{ip: addr, port: featureRegister.Port, commands: featureRegister.Commands}

	services = append(services, service)
	message := fmt.Sprintf("\n\nServiço adicionado. IP/Porta: %s Prefixos: %s", service.getIpAndPort(), service.commands)
	fmt.Println(message)
	return 200, message
}

func listenToClient() {
	router := http.NewServeMux()

	router.HandleFunc("POST /command", func(w http.ResponseWriter, r *http.Request) {
		var command common.Command
		_ = json.NewDecoder(r.Body).Decode(&command)
		statusCode, message := handleClientCommand(command)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(message))
	})

	http.ListenAndServe(":"+strconv.Itoa(ClientPort), c.Handler(router))
}

func handleClientCommand(command common.Command) (int, string) {
	commandStr := strings.TrimSpace(command.Command)

	fmt.Println("RECEBEU O COMANDO: " + commandStr)

	var endpoint = ""
	var index = slices.IndexFunc(services, func(s Service) bool {
		var ok = false
		endpoint, ok = s.commands[commandStr]
		return ok
	})

	if index == -1 {
		return 500, "NÃO HÁ SERVIÇO CAPAZ DE RESPONDER SUA SOLICITAÇÃO"
	}

	service := services[index]
	url := fmt.Sprintf("http://%s:%d/%s", service.ip, service.port, endpoint)

	body, err := json.Marshal(command.Arguments)

	if err != nil {
		return 400, "FALHA AO DECODIFICAR COMANDO"
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	defer resp.Body.Close()

	if err != nil {
		return 500, "NÃO FOI POSSÍVEL EXECUTAR O COMANDO"
	}

	result, _ := io.ReadAll(resp.Body)
	fmt.Printf("RESPONDEU: %s", result)

	return 200, string(result)
}
