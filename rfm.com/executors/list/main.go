package main

import (
	"bytes"
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "os"
	"rfm.com/common"
	"rfm.com/executors/list/api"
	_ "rfm.com/executors/list/docs"
	_ "rfm.com/executors/list/model"
	_ "slices"
	"strconv"
)

const DiscoveryPort = 7070

var port string = "8888"
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

// @host		localhost:8888
// @BasePath	/
func main() {
	go communicateDiscovery()
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("GET /list", api.GetDirectories)
	http.HandleFunc("POST /list", api.GetDirectoriesByBody)
	_ = http.ListenAndServe(":"+port, nil)
}

func communicateDiscovery() {
	for {

		selfPort, _ := strconv.Atoi(port)
		featureRegister := common.FeatureRegister{Port: selfPort, Commands: map[string]string{"list": "/list"}}
		data, _ := json.Marshal(featureRegister)

		resp, err := http.Post("http://"+discoveryIP+":"+strconv.Itoa(DiscoveryPort)+"/register", "application/json", bytes.NewBuffer(data))

		if err == nil && resp.StatusCode == 200 {
			return
		}
	}
}
