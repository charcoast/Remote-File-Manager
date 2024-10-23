package main

import (
	"encoding/json"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	_ "os"
	"rfm.com/commom"
	"rfm.com/executors/list/api"
	_ "rfm.com/executors/list/docs"
	_ "rfm.com/executors/list/model"
	_ "slices"
	"strconv"
)

const DiscoveryPort = 7070

var port string
var discoveryIP string
var prefixes = []string{"list", "li", "ls"}

//	@title			Remote-File-Manager - List Executor
//	@version		1.0
//	@description	This is one of the executors that works in the Remote-File-Manager application.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func main() {
	go communicateDiscovery()
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("GET /list", api.GetDirectories)
	http.HandleFunc("POST /list", api.GetDirectoriesByBody)
	_ = http.ListenAndServe(":8080", nil)
}

func communicateDiscovery() {
	for {
		conn, err := net.Dial("tcp", discoveryIP+":"+strconv.Itoa(DiscoveryPort))
		if err != nil {
			continue
		}

		selfPort, _ := strconv.Atoi(port)
		featureRegister := commom.FeatureRegister{Port: selfPort, Prefixes: prefixes}
		data, _ := json.Marshal(featureRegister)
		sendData(conn, data)
		closeConnection(conn)
		return
	}
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func sendData(conn net.Conn, b []byte) {
	_, err := fmt.Fprintf(conn, string(b)+"\n")
	if err != nil {
		sysOut(err.Error())
		return
	}
}

func sysOut(value interface{}) {
	fmt.Println(value)
}
