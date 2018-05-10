package node

import (
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Test - random test function
func Test() {
}

// CreateNode - create a Node contract
func CreateNode() (string, error) {
	url := "http://localhost:3001/api/node/create"

	res, err := utils.SendRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.CreateNode", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.CreateNode", err)
	}

	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}

// GetNodeAddress - get node address from owner lookup
func GetNodeAddress() (string, error) {
	url := "http://localhost:3001/api/node/"

	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.GetNodeAddress", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.GetNodeAddress", err)
	}

	address := api.Response.(string)

	return address, nil //node address
}

// SetNodeData - set data for a Node contract
func SetNodeData(nodeAddress string, data map[string]interface{}) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/%s/data", nodeAddress)

	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", fmt.Errorf("%v/node.SetNodeData", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.SetNodeData", err)
	}

	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}

// ApplyToPool - apply to a pool [Need to implement new API]
func ApplyToPool(nodeAddress, poolAddress string) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/%s/apply/%s", nodeAddress, poolAddress)

	res, err := utils.SendRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.ApplyToPool", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.CreateNode", err)
	}

	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}

// CheckPoolApplication - check the status of your pool application [Need to implement new API]
func CheckPoolApplication(nodeAddress, poolAddress string) (string, error) {
	url := fmt.Sprintf("http://localhost:3001/api/node/%s/application/%s", nodeAddress, poolAddress)

	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/node.CheckPoolApplication", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/node.CheckPoolApplication", err)
	}

	response := api.Response.(map[string]interface{})
	status := response["status"].(string)

	return status, nil // pool status
}

// StartEdgeNode - start edge node server
func StartEdgeNode() (string, error) {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	err := clientHTTP.Call("GladiusEdge.Start", nil, &reply)
	if err != nil {
		return "", fmt.Errorf("%v/node.StopEdgeNode", err)
	}
	return reply, nil
}

// StopEdgeNode - stop edge node server
func StopEdgeNode() (string, error) {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	err := clientHTTP.Call("GladiusEdge.Stop", nil, &reply)
	if err != nil {
		return "", fmt.Errorf("%v/node.StopEdgeNode", err)
	}

	return reply, nil
}

// StatusEdgeNode - status of edge node server
func StatusEdgeNode() (string, error) {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	err := clientHTTP.Call("GladiusEdge.Status", nil, &reply)
	if err != nil {
		return "", fmt.Errorf("%v/node.StatusEdgeNode", err)
	}

	return reply, nil
}
