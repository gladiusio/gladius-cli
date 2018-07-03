package node

import (
	"fmt"
	"os"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/powerman/rpc-codec/jsonrpc2"
	log "github.com/sirupsen/logrus"
)

// LogFile - Where the logs are stored
var LogFile *os.File

// Test - random test function
func Test() {
}

// CreateNode - create a Node contract
func CreateNode() (string, error) {
	url := "http://localhost:3001/api/node/create"

	log.WithFields(log.Fields{"file": "node.go", "func": "CreateNode"}).Debug("POST: ", url)
	// use the custom sendRequest to send something to the control daemon api
	res, err := utils.SendRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.CreateNode", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "CreateNode"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.CreateNode", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "CreateNode"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}

// GetNodeAddress - get node address from owner lookup
func GetNodeAddress() (string, error) {
	url := "http://localhost:3001/api/node/"

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeAddress"}).Debug("GET: ", url)
	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.GetNodeAddress", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeAddress"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.GetNodeAddress", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeAddress"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	address := response["address"].(string)

	return address, nil //node address
}

// GetNodeData - get node address from owner lookup
func GetNodeData(address string) (map[string]interface{}, error) {
	url := "http://localhost:3001/api/node/" + address + "/data"

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeData"}).Debug("GET: ", url)
	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%v/node.GetNodeAddress", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeData"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return nil, fmt.Errorf("%v/node.GetNodeAddress", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "GetNodeData"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})

	return response, nil //node data
}

// SetNodeData - set data for a Node contract
func SetNodeData(nodeAddress string, data map[string]interface{}) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/%s/data", nodeAddress)

	log.WithFields(log.Fields{"file": "node.go", "func": "SetNodeData"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", fmt.Errorf("%v/node.SetNodeData", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "SetNodeData"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.SetNodeData", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "SetNodeData"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}

// ApplyToPool - apply to a pool [Need to implement new API]
func ApplyToPool(nodeAddress, poolAddress string) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/%s/apply/%s", nodeAddress, poolAddress)

	log.WithFields(log.Fields{"file": "node.go", "func": "ApplyToPool"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.ApplyToPool", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "ApplyToPool"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.CreateNode", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "ApplyToPool"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}

// CheckPoolApplication - check the status of your pool application [Need to implement new API]
func CheckPoolApplication(nodeAddress, poolAddress string) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/%s/application/%s", nodeAddress, poolAddress)

	log.WithFields(log.Fields{"file": "node.go", "func": "CheckPoolApplication"}).Debug("GET: ", url)
	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.CheckPoolApplication", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "CheckPoolApplication"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.CheckPoolApplication", err)
	}

	log.WithFields(log.Fields{"file": "node.go", "func": "CheckPoolApplication"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	status := response["status"].(string)
	return status, nil // pool status
}

// StartNetworkNode - start networking node server
func StartNetworkNode() (string, error) {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	err := clientHTTP.Call("GladiusEdge.Start", nil, &reply)
	if err != nil {
		return "", fmt.Errorf("%v/node.StopNetworkNode", err)
	}
	return reply, nil
}

// StopNetworkNode - stop network node server
func StopNetworkNode() (string, error) {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	err := clientHTTP.Call("GladiusEdge.Stop", nil, &reply)
	if err != nil {
		return "", fmt.Errorf("%v/node.StopNetworkNode", err)
	}

	return reply, nil
}

// StatusNetworkNode - status of network node server
func StatusNetworkNode() (string, error) {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	err := clientHTTP.Call("GladiusEdge.Status", nil, &reply)
	if err != nil {
		return "", fmt.Errorf("%v/node.StatusNetworkNode", err)
	}

	return reply, nil
}

func init() {
	// set up the logger
	switch utils.LogLevel {
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.FatalLevel)
	}

	LogFile, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Warning("Failed to log to file, using default stderr")
	}

	log.SetOutput(LogFile)
}
