package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gladiusio/gladius-cli/internal"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Node - properties of a node
type Node struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Data    struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		IPAddress string `json:"ipAddress"`
		Status    string `json:"status"`
	} `json:"data"`
}

// ApiResponse - standard response from the control daemon api
type ApiResponse struct {
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
	Endpoint string      `json:"endpoint"`
}

// For control over HTTP client headers,
// redirect policy, and other settings,
// create an HTTP client
var client = &http.Client{
	Timeout: time.Second * 10, //10 second timeout
}

// Test ...
func Test(myNode Node) {
	test := myNode.Data
	fmt.Println((test))

	_, err := json.Marshal(test)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// PostSettings - posts user settings to the api [Deprecated in CDv2]
// func PostSettings(filename string) bool {
// 	url := "http://localhost:3000/api/settings/start"
//
// 	envFile, err := utils.GetEnvMap("env.toml")
//
// 	envData := envFile["environment"] // only use what's in the environment section
//
// 	_, err = utils.SendRequest(client, "POST", url, envData)
// 	if err != nil {
// 		log.Fatal("POST-postSettings(): ", err)
// 		return false
// 	}
//
// 	return true
// }

// GetSettings - get settings from API [Needs to be implemented in CDv2]
func GetSettings() {
	url := "http://localhost:3000/api/settings/"

	res, err := utils.SendRequest(client, "GET", url, nil)
	if err != nil {
		log.Fatal("GET-getSettings(): ", err)
	}

	fmt.Println(res)
}

// CreateNode - create a Node contract
func CreateNode() (string, error) {
	url := "http://localhost:3001/api/node/create"

	res, err := utils.SendRequest(client, "POST", url, nil)
	if err != nil {
		return "", err
	}

	api, err := ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", err
	}

	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	fmt.Println("address: ", txHash["value"])

	return txHash["value"].(string), nil //tx hash
}

// GetNodeAddress - get node address from owner lookup
func GetNodeAddress() string {
	url := "http://localhost:3000/api/node"

	res, err := utils.SendRequest(client, "GET", url, nil)
	if err != nil {
		log.Fatal("GET-getNodeAddress(): ", err)
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(res), &data)

	return data["address"].(string) // node address
}

// SetNodeData - set data for a Node contract
func SetNodeData(nodeAddress string, myNode Node) (string, error) {
	url := fmt.Sprintf("http://localhost:3000/api/node/%s/data", nodeAddress)

	res, err := utils.SendRequest(client, "POST", url, myNode.Data)
	if err != nil {
		log.Fatal("POST-setNodeData(): ", err)
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(res), &data)

	if data["txHash"] == nil {
		return "", errors.New("ERROR CREATING NODE")
	}

	return data["txHash"].(string), nil // tx hash
}

// ApplyToPool - apply to a pool
func ApplyToPool(nodeAddress, poolAddress string) (string, error) {
	url := fmt.Sprintf("http://localhost:3000/api/node/%s/apply/%s", nodeAddress, poolAddress)

	res, err := utils.SendRequest(client, "POST", url, nil)
	if err != nil {
		log.Fatal("POST-postSettings(): ", err)
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(res), &data)

	if data["tx"] == nil {
		return "", errors.New("ERROR APPLYING TO POOL: " + res)
	}

	return data["tx"].(string), nil // tx hash
}

// CheckPoolApplication - check the status of your pool application
func CheckPoolApplication(nodeAddress, poolAddress string) string {
	url := fmt.Sprintf("http://localhost:3000/api/node/%s/status/%s", nodeAddress, poolAddress)

	res, err := utils.SendRequest(client, "GET", url, nil)
	if err != nil {
		log.Fatal("GET-getPoolStatus(): ", err)
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(res), &data)

	return data["status"].(string) // application status
}

// CheckTx - check status of tx hash
func CheckTx(tx string) (bool, error) {
	url := fmt.Sprintf("http://localhost:3001/api/status/tx/%s", tx)

	res, err := utils.SendRequest(client, "GET", url, nil)
	if err != nil {
		log.Fatal("POST-checkTx(): ", err)
		return false, err
	}

	api, err := ControlDaemonHandler([]byte(res))
	if err != nil {
		return false, err
	}

	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	fmt.Println("address: ", txHash["value"])

	return txHash["complete"], nil // tx completion status
}

// WaitForTx - wait for the tx to complete
func WaitForTx(tx string) bool {
	status := CheckTx(tx)

	for status == false {
		status = CheckTx(tx)
		fmt.Printf("Tx: %s\t Status: Pending\r", tx)
	}

	fmt.Printf("\nTx: %s\t Status: Successful\n", tx)
	return true
}

// Should add errors for the edge node functions below

// StartEdgeNode - start edge node server
func StartEdgeNode() string {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	clientHTTP.Call("GladiusEdge.Start", nil, &reply)

	return reply
}

// StopEdgeNode - stop edge node server
func StopEdgeNode() string {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	clientHTTP.Call("GladiusEdge.Stop", nil, &reply)

	return reply
}

// StatusEdgeNode - status of edge node server
func StatusEdgeNode() string {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://localhost:5000/rpc")
	defer clientHTTP.Close()

	var reply string

	// Synchronous call using positional params and TCP.
	clientHTTP.Call("GladiusEdge.Status", nil, &reply)

	return reply
}

// handle the control daemon responses
func ControlDaemonHandler(_res []byte) (ApiResponse, error) {
	var response = ApiResponse{}

	err := json.Unmarshal(_res, &response)
	if err != nil {
		return ApiResponse{}, err
	}

	if !response.Success {
		return ApiResponse{}, errors.New("API ERROR: " + response.Message)
	}

	return response, nil
}
