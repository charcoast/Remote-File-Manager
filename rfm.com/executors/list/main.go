package main

import (
	"fmt"
	"net/http"
	_ "os"
	_ "slices"
	"strconv"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"rfm.com/common"
	"rfm.com/executors/list/api"
	_ "rfm.com/executors/list/docs"
	_ "rfm.com/executors/list/model"
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
	commands := map[string]string{"ls": "/list"}
	featureRegister := common.FeatureRegister{Port: SelfPort, Commands: commands}
	common.CommunicateDiscovery(featureRegister)
}

// func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		// Handle preflight requests
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	}
// }

func sysOut(value interface{}) {
	fmt.Println(value)
}
