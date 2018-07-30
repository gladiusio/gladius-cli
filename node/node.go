package node

import (
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/powerman/rpc-codec/jsonrpc2"
	log "github.com/sirupsen/logrus"
)

// GetApplication- get node application from pool
func GetApplication(address string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/applications/%s/view", address)

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

	valid := application["Valid"].(bool)
	accepted := application["Bool"].(bool)

	if valid {
		if accepted {
			return "Accepted", nil
		}
	} else {
		return "Pending", nil
	}

	return "Rejected", nil
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
		return "", utils.HandleError(err, "Error starting the node networking daemon. Make sure it's running!", "node.StartNetworkNode")
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
		return "", utils.HandleError(err, "Error stopping the node networking daemon. Make sure it's running!", "node.StopNetworkNode")
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
		return "", utils.HandleError(err, "Error communicating with the node networking daemon. Make sure it's running!", "node.StatusNetworkNode")
	}

	return reply, nil
}
