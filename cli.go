package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type node struct {
	Address string `json:"address"`
}

// For control over HTTP client headers,
// redirect policy, and other settings,
// create a Client
// A Client is an HTTP client
var client = &http.Client{
	Timeout: time.Second * 10, //10 second timeout
}

func main() {

	controlSettings()

}

func controlSettings() {
	url := "http://localhost:3000/api/node"

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	req.Header.Set("User-Agent", "gladius-cli")

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// myNode := node{}
	// jsonErr := json.Unmarshal(body, &myNode)
	// if jsonErr != nil {
	// 	log.Fatal(jsonErr)
	// }

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer res.Body.Close()

	fmt.Println(string(body))
}
