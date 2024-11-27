package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"github.com/thoas/go-funk"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net"
	"net/http"
	"rfm.com/common"
	"rfm.com/discovery/model"
	"slices"
	"strconv"
	"strings"
)

const DiscoveryPort = 7070
const ClientPort = 9090

var services []Service
var c *cors.Cors
var db *gorm.DB = nil

func main() {

	initDB()

	fmt.Println("****RUNNING****")
	c = cors.New(cors.Options{AllowedOrigins: []string{"http://167.234.232.150"},
		AllowCredentials: true,
		Debug:            true})

	go listenToServices()
	listenToClient()
}

func initDB() {
	var err error

	db, err = gorm.Open(sqlite.Open("discovery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	var user = model.User{}
	db.First(&user, "username = ?", "admin")

	if user.Username != "" {
		return
	}

	hashPass, err := hashPassword("2M5w93ASyNFh=O>liMYU7")
	db.Create(&model.User{Username: "admin", Password: hashPass})
}

func listenToServices() {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", handleRegisterService)
	err := http.ListenAndServe(":"+strconv.Itoa(DiscoveryPort), c.Handler(router))
	if err != nil {
		panic(err)
	}
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

func handleAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, plainPass, _ := r.BasicAuth()

		var user model.User
		db.First(&user, "username = ?", username)

		if user.Username == "" {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(401)
			return
		}

		valid := VerifyPassword(plainPass, user.Password)

		if !valid {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(401)
			return
		}

		h.ServeHTTP(w, r)

	})

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func listenToClient() {
	router := http.NewServeMux()

	router.HandleFunc("POST /command", func(w http.ResponseWriter, r *http.Request) {
		var command common.Command
		_ = json.NewDecoder(r.Body).Decode(&command)
		statusCode, header, message := handleClientCommand(command)

		if header != nil {
			w.Header().Set("Content-Type", header.Get("Content-Type"))
			w.Header().Set("Content-Disposition", header.Get("Content-Disposition"))
			w.Header().Set("Content-Length", header.Get("Content-Length"))
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(statusCode)
		w.Write(message)
	})

	http.ListenAndServe(":"+strconv.Itoa(ClientPort), c.Handler(handleAuthentication(router)))
}

func handleClientCommand(command common.Command) (int, http.Header, []byte) {
	commandStr := strings.TrimSpace(command.Command)

	fmt.Println("RECEBEU O COMANDO: " + commandStr)

	var endpoint = ""
	var index = slices.IndexFunc(services, func(s Service) bool {
		var ok = false
		endpoint, ok = s.commands[commandStr]
		return ok
	})

	if index == -1 {
		return 500, nil, []byte("NÃO HÁ SERVIÇO CAPAZ DE RESPONDER SUA SOLICITAÇÃO")
	}

	service := services[index]
	url := fmt.Sprintf("http://%s:%d/%s", service.ip, service.port, endpoint)

	body, err := json.Marshal(command.Arguments)

	if err != nil {
		return 400, nil, []byte("FALHA AO DECODIFICAR COMANDO")
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	defer resp.Body.Close()

	if err != nil {
		return 500, nil, []byte("NÃO FOI POSSÍVEL EXECUTAR O COMANDO")
	}

	result, _ := io.ReadAll(resp.Body)
	fmt.Printf("RESPONDEU: %s", result)

	return 200, resp.Header, result
}
