package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"rfm.com/common"
	_ "rfm.com/discovery/docs"
	"slices"
	"strconv"
)

const DiscoveryPort = 7070
const ClientPort = 9090

var services []Service

//	@title			Remote-File-Manager - Discovery
//	@version		1.0
//	@description	This is the discovery service that receive connections from executor and the client
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:7070
// @BasePath	/
func main() {

	fmt.Println("****RUNNING****")

	go listenToServices()
	listenToClient()

}

func listenToServices() {

	router := http.NewServeMux()

	router.HandleFunc("POST /register", handleRegisterService)
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.ListenAndServe(":"+strconv.Itoa(DiscoveryPort), router)

}

// handleRegisterService
// @Summary
// @Description Register a service in the service discovery
// @ID register-service
// @Accept json
// @Produce text/plain
// @Param request body common.FeatureRegister true "query params"
// @Success 200 {string} []string "ok"
// @Failure 500 {object} string "Can not get the directories"
// @Router /register [post]
func handleRegisterService(w http.ResponseWriter, r *http.Request) {
	var featureRegister common.FeatureRegister
	_ = json.NewDecoder(r.Body).Decode(&featureRegister)
	statusCode, message := handleServiceDiscovery(r.RemoteAddr, featureRegister)

	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func handleServiceDiscovery(addr string, featureRegister common.FeatureRegister) (int, string) {

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
	body, err := json.Marshal(common.Command{Command: command})

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
