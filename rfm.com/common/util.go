package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func CommunicateDiscovery(featureRegister FeatureRegister) {
	for {
		url := fmt.Sprintf("http://%s:%d/register", DiscoveryDomain, DiscoveryPort)
		jsonValue, _ := json.Marshal(featureRegister)
		request := bytes.NewBuffer(jsonValue)
		response, _ := http.Post(url, "application/json", request)
		if response != nil && response.StatusCode == 200 {
			return
		}
	}
}
