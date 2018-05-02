package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BurntSushi/toml"
)

// custom function to make sending request less of a pain in the arse
func SendRequest(client *http.Client, requestType, url string, data interface{}) (string, error) {

	b := bytes.Buffer{}

	if data != nil {
		jsonPayload, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		b = *bytes.NewBuffer(jsonPayload)
	}

	// Build the request
	req, err := http.NewRequest(requestType, url, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "gladius-cli")
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", err
	}

	// Defer the closing of the body
	defer res.Body.Close()

	return string(body), nil //tx
}

func GetEnvMap(filename string) (map[string]map[string]string, error) {
	b, err := ioutil.ReadFile(filename) // read env file
	if err != nil {
		fmt.Println("Error reading: " + filename)
		return nil, err
	}

	var envFile = make(map[string]map[string]string)

	if _, err := toml.Decode(string(b), &envFile); err != nil { // turn file into mapping
		fmt.Println("Error decoding")
		return nil, err
	}

	return envFile, nil
}

func WriteToEnv(section, key, value, source, destination string) error {
	b, err := ioutil.ReadFile(source) // read env file
	if err != nil {
		fmt.Println("Error reading: " + source)
		return err
	}

	var envFile = make(map[string]map[string]string)

	if _, err := toml.Decode(string(b), &envFile); err != nil { // turn file into mapping
		fmt.Println("Error decoding")
	}

	envFile[section][key] = value

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(envFile); err != nil {
		fmt.Println("Error encoding")
		return err
	}

	if err = ioutil.WriteFile(destination, (*buf).Bytes(), 0644); err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}
