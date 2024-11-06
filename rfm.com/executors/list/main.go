package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
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

var SelfPort = common.GetRandomPort()

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
	sysOut(SelfPort)
	go communicateDiscovery()
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("GET /list", api.GetDirectories)
	http.HandleFunc("POST /list", api.GetDirectoriesByBody)
	cors.AllowAll()
	_ = http.ListenAndServe(":"+strconv.Itoa(SelfPort), nil)
}

func communicateDiscovery() {
	for {
		url := fmt.Sprintf("http://%s:%d/register", common.DiscoveryDomain, common.DiscoveryPort)
		commands := map[string]string{"ls": "/list"}
		featureRegister := common.FeatureRegister{Port: SelfPort, Commands: commands}
		jsonValue, _ := json.Marshal(featureRegister)
		request := bytes.NewBuffer(jsonValue)
		response, _ := http.Post(url, "application/json", request)
		if response != nil && response.StatusCode == 200 {
			return
		}
	}
}

func sysOut(value interface{}) {
	fmt.Println(value)
}
