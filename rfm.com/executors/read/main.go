package main

import (
	"fmt"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "os"
	"rfm.com/common"
	"rfm.com/executors/read/api"
	_ "rfm.com/executors/read/docs"
	_ "rfm.com/executors/read/model"
	_ "slices"
	"strconv"
)

var SelfPort = common.GetRandomPort()

//	@title			Remote-File-Manager - Read Executor
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
	http.HandleFunc("POST /read", api.ReadFile)
	cors.AllowAll()
	_ = http.ListenAndServe(":"+strconv.Itoa(SelfPort), nil)
}

func communicateDiscovery() {
	commands := map[string]string{"ls": "/list"}
	featureRegister := common.FeatureRegister{Port: SelfPort, Commands: commands}
	common.CommunicateDiscovery(featureRegister)
}

func sysOut(value interface{}) {
	fmt.Println(value)
}
