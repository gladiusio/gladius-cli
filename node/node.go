package node

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetApplication - get node application from pool
func GetApplication(poolAddress string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://localhost:%d/api/node/applications/%s/view", viper.GetInt("Ports.NetworkGateway"), poolAddress)

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
	url := fmt.Sprintf("http://localhost:%d/api/node/applications/%s/new", viper.GetInt("Ports.NetworkGateway"), poolAddress)

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

// Start - start network gateway and edged
func Start() (string, error) {
	timeoutURL := fmt.Sprintf("http://localhost:%d/service/set_timeout", viper.GetInt("Ports.Guardian"))
	startURL := fmt.Sprintf("http://localhost:%d/service/set_state/all", viper.GetInt("Ports.Guardian"))

	timeout := make(map[string]int)
	timeout["timeout"] = 3

	running := make(map[string]bool)
	running["running"] = true

	log.WithFields(log.Fields{"file": "node.go", "func": "Start"}).Debug("POST: ", timeoutURL)
	_, err := utils.SendRequest("POST", timeoutURL, timeout)
	if err != nil {
		return "Failed to set timeout", utils.HandleError(err, "", "node.Start")
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "Start"}).Debug("POST: ", startURL)
	_, err = utils.SendRequest("PUT", startURL, running)
	if err != nil {
		return "Failed to star one or more modules", utils.HandleError(err, "", "node.Start")
	}

	return "Started modules", nil
}

// Stop - stop network gateway and edged
func Stop() (string, error) {
	stopURL := fmt.Sprintf("http://localhost:%d/service/set_state/all", viper.GetInt("Ports.Guardian"))

	running := make(map[string]bool)
	running["running"] = false

	log.WithFields(log.Fields{"file": "node.go", "func": "Start"}).Debug("POST: ", stopURL)
	_, err := utils.SendRequest("PUT", stopURL, running)
	if err != nil {
		return "Failed to stop one or both modules", utils.HandleError(err, "", "node.Start")
	}

	return "Stopped modules", nil
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

func NeedUpdate() (bool, error) {
	// get the official versions
	res, err := utils.SendRequest("GET", "https://gladius-version.nyc3.digitaloceanspaces.com/version.json", nil)
	if err != nil {
		return false, err
	}

	// parse response
	var response = make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		return false, err
	}

	officialVersions := response

	// get the current versions
	currentVersion := make(map[string]string)
	currentVersion["gladius-guardian"], _ = GetVersion("guardian")
	currentVersion["gladius-edged"], _ = GetVersion("edged")
	currentVersion["gladius-network-gateway"], _ = GetVersion("network-gateway")

	var needUpdate [3]bool

	// compare
	count := 0
	for ver := range officialVersions {
		if officialVersions[ver] != currentVersion[ver] {
			needUpdate[count] = false
		} else {
			needUpdate[count] = true
		}
		count++
	}

	// if anything here is false then we need an update
	if needUpdate[0] && needUpdate[1] && needUpdate[2] {
		return false, nil
	} else {
		return true, utils.HandleError(errors.New("One or more of your modules is out of date"), "One or more of your modules is out of date", "node.needUpdate")
	}
}
