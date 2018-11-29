package node

import (
	"encoding/json"
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetApplication - get node application from pool
func GetApplication(poolAddress string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/applications/%s/view", poolAddress)

	log.WithFields(log.Fields{"file": "node.go", "func": "GetApplication"}).Debug("GET: ", url)
	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return nil, utils.HandleError(err, "", "node.GetNodeData")
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "GetApplication"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return nil, utils.HandleError(err, "", "node.GetNodeData")
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeData"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	data := response["profile"].(map[string]interface{})

	return data, nil //node data
}

// ApplyToPool - apply to a pool
func ApplyToPool(poolAddress string, data map[string]interface{}) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/applications/%s/new", poolAddress)

	log.WithFields(log.Fields{"file": "node.go", "func": "ApplyToPool"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", utils.HandleError(err, "", "node.AppyToPool")
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "ApplyToPool"}).Debug("Response recieved, piping through the response handler")
	_, err = utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", utils.HandleError(err, "", "node.AppyToPool")
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "ApplyToPool"}).Debug("Decoding response fields")

	return "success", nil //tx hash
}

// CheckPoolApplication - check the status of your pool application
func CheckPoolApplication(poolAddress string) (string, error) {
	application, err := GetApplication(poolAddress)
	if err != nil {
		return "", utils.HandleError(err, "", "node.CheckPoolApplication")
	}

	if application == nil {
		return "No Application Found", nil
	}

	pending := application["pending"].(bool)
	accepted := application["approved"].(bool)

	if pending {
		return "Pending", nil
	}

	if accepted {
		return "Accepted", nil
	}

	return "Rejected", nil
}

// StatusNetworkNode - status of network node server
func StatusNetworkNode() (string, error) {
	url := "http://localhost:8080"

	log.WithFields(log.Fields{"file": "node.go", "func": "StatusNetworkNode"}).Debug("GET: ", url)
	_, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "Offline", utils.HandleError(err, "", "node.StatusNetworkNode")
	}

	return "Online", nil
}

// GetVersion - get individual version number from module
func GetVersion(module string) (string, error) {
	var port int
	switch module {
	case "guardian":
		port = viper.GetInt("Ports.Guardian")
	case "edged":
		port = viper.GetInt("Ports.EdgeD")
	case "network-gateway":
		port = viper.GetInt("Ports.NetworkGateway")
	default:
		port = 0
	}

	if port == 0 {
		return "", fmt.Errorf("Module %s not found", module)
	}
	res, err := utils.SendRequest("GET", fmt.Sprintf("http://localhost:%d/version", port), nil)
	if err != nil {
		return "", err
	}

	var response = make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		return "", err
	}

	res1 := response["response"].(map[string]interface{})
	version := res1["version"].(string)

	return version, nil
}
