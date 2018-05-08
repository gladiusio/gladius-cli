package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

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

// SendRequest - custom function to make sending request less of a pain in the arse
func SendRequest(requestType, url string, data interface{}) (string, error) {

	b := bytes.Buffer{}

	// if data present, turn it into a bytesBuffer(jsonPayload)
	if data != nil {
		jsonPayload, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("%v:json.Marshall/utils.sendRequest", err)
		}
		b = *bytes.NewBuffer(jsonPayload)
	}

	// Build the request
	req, err := http.NewRequest(requestType, url, &b)
	if err != nil {
		return "", fmt.Errorf("%v:http.NewRequest/utils.SendRequest", err)
	}

	req.Header.Set("User-Agent", "gladius-cli")
	req.Header.Set("Content-Type", "application/json")

	// add the X-Auth header and pass in a user password if requestType is anything but GET
	if strings.Compare(requestType, "GET") != 0 {
		password := ""
		prompt := &survey.Password{
			Message: "Please type your password: ",
		}
		survey.AskOne(prompt, &password, nil)

		req.Header.Set("X-Authorization", password)
	}

	// Send the request via a client
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%v:client.Do/utils.SendRequest", err)
	}

	// read the body of the response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", fmt.Errorf("%v:ioutil.ReadAll/utils.SendRequest", err)
	}

	// Defer the closing of the body
	defer res.Body.Close()

	return string(body), nil //tx
}

// CheckTx - check status of tx hash
func CheckTx(tx string) (bool, error) {
	url := fmt.Sprintf("http://localhost:3001/api/status/tx/%s", tx)

	res, err := SendRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("%v/utils.CheckTx", err)
	}

	api, err := ControlDaemonHandler([]byte(res))
	if err != nil {
		return false, fmt.Errorf("%v/utils.CheckTx", err)
	}

	response := api.Response.(map[string]interface{})

	if response["complete"] == false {
		return false, nil
	}

	return response["complete"].(bool), nil // tx completion status
}

// WaitForTx - wait for the tx to complete
func WaitForTx(tx string) (bool, error) {
	status, err := CheckTx(tx)
	if err != nil {
		return false, fmt.Errorf("%v/utils.WaitForTx", err)
	}

	for !status {
		status, err = CheckTx(tx)
		if err != nil {
			time.Sleep(1 * time.Second)
			return false, fmt.Errorf("%v/utils.WaitForTx", err)
		}
		fmt.Printf("Tx: %s\t Status: Pending\r", tx)
	}

	fmt.Printf("\nTx: %s\t Status: Successful\n", tx)
	return true, nil
}

// ControlDaemonHandler - handler for the API responses
func ControlDaemonHandler(_res []byte) (ApiResponse, error) {
	var response = ApiResponse{}

	err := json.Unmarshal(_res, &response)
	if err != nil {
		return ApiResponse{}, fmt.Errorf("%v:json.Unmarshall/utils.ControlDaemonHandler", err)
	}

	if !response.Success {
		return ApiResponse{}, fmt.Errorf("%s:utils.ControlDaemonHandler", response.Message)
	}

	return response, nil
}

// GetIP - Retrieve the current machine's external IP address
func GetIP() (string, error) {
	res, err := SendRequest("GET", "http://ipv4.myexternalip.com/raw", nil)
	if err != nil {
		return "", fmt.Errorf("%v:utils.GetIP", err)
	}
	return res, nil
}

// ############ DEPRECATED ############

// custom function to return a mapping of the environment file (has to be .toml)
// this technically works but reading from *.toml is deprecated
func GetEnvMap(filename string) (map[string]map[string]string, error) {
	// read env file
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading: " + filename)
		return nil, err
	}

	// decode the file and put it into envFile
	var envFile = make(map[string]map[string]string)
	if _, err := toml.Decode(string(b), &envFile); err != nil {
		fmt.Println("Error decoding")
		return nil, err
	}

	return envFile, nil
}

// custom function to return a mapping of the environment file (has to be .toml)
// this technically works but writing to *.toml is deprecated
func WriteToEnv(section, key, value, source, destination string) error {
	// read the file
	b, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Println("Error reading: " + source)
		return err
	}

	// decode and put it into the mapping
	var envFile = make(map[string]map[string]string)
	if _, err := toml.Decode(string(b), &envFile); err != nil {
		fmt.Println("Error decoding")
	}

	// add a new {key : value} pair
	envFile[section][key] = value

	// re-encode the mapping
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(envFile); err != nil {
		fmt.Println("Error encoding")
		return err
	}

	// re-write the file
	if err = ioutil.WriteFile(destination, (*buf).Bytes(), 0644); err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}
