package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	log "github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// APIResponse - standard response from the control daemon api
type APIResponse struct {
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
	TxHash   interface{} `json:"txHash"`
	Endpoint string      `json:"endpoint"`
}

// ErrorResponse - custom error struct
type ErrorResponse struct {
	UserMessage string
	LogError    string
	Path        string
}

var cachedPassphrase string

// RequestTimeout - Request timeout in seconds
var RequestTimeout int

// Error - for the dev/logger
func (e *ErrorResponse) Error() string {
	return e.LogError
}

// Message - for the user
func (e *ErrorResponse) Message() string {
	return e.UserMessage
}

// For control over HTTP client headers,
// redirect policy, and other settings,
// create an HTTP client
var client = &http.Client{
	Timeout: time.Second * time.Duration(RequestTimeout),
}

// SendRequest - custom function to make sending api requests less of a pain
// in the arse.
func SendRequest(requestType, url string, data interface{}) (string, error) {

	b := bytes.Buffer{}

	// if data present, turn it into a bytesBuffer(jsonPayload)
	if data != nil {
		jsonPayload, err := json.Marshal(data)
		if err != nil {
			return "", HandleError(err, "Invalid Data", ":json.Marshall/SendRequest")
		}
		b = *bytes.NewBuffer(jsonPayload)
	}

	// Build the request
	req, err := http.NewRequest(requestType, url, &b)
	if err != nil {
		return "", HandleError(err, "Could not build request", ":http.NewRequest/SendRequest")
	}

	req.Header.Set("User-Agent", "gladius-cli")
	req.Header.Set("Content-Type", "application/json")

	// if you're writing then ask for password
	if strings.Compare(requestType, "GET") != 0 {
		// if the password is cached then use it
		if strings.Compare(cachedPassphrase, "") != 0 {
			req.Header.Set("X-Authorization", cachedPassphrase)
		} else { //else ask
			password := AskPassphrase()
			cachedPassphrase = password
			req.Header.Set("X-Authorization", password)
		}
	}

	// Send the request via a client
	res, err := client.Do(req)
	if err != nil {
		return "", HandleError(err, "Could not send request", ":client.Do/SendRequest")
	}

	// read the body of the response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", HandleError(err, "Could not build request", ":ioutil.ReadAll/SendRequest")
	}

	// Defer the closing of the body
	defer res.Body.Close()

	return string(body), nil //tx
}

// CheckTx - check status of tx.
// Perform a single check on a tx.
func CheckTx(tx string) (bool, error) {
	url := fmt.Sprintf("http://localhost:3001/api/status/tx/%s", tx)

	res, err := SendRequest("GET", url, nil)
	if err != nil {
		return false, HandleError(err, "", "utils.CheckTx")
	}

	api, err := ControlDaemonHandler([]byte(res))
	if err != nil {
		return false, HandleError(err, "", "utils.CheckTx")
	}

	response := api.Response.(map[string]interface{})

	if response["complete"] == false {
		return false, nil
	}

	return response["complete"].(bool), nil // tx completion status
}

// WaitForTx - wait for a tx on the blockchain to complete.
// Queries the API every second to see if tx is complete.
func WaitForTx(tx string) (bool, error) {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan error) // this is the exit condition channel

	println()

	// hit the status API every 1 second
	go func() {
		count := 0
		for {
			select {
			case <-ticker.C:
				status, err := CheckTx(tx)
				if err != nil {
					quit <- err // if there's an error here then pump it into the channel
				}
				if status {
					quit <- nil // if the tx went through, pump a nil error into the channel
				}
				switch count {
				case 0:
					fmt.Printf("Tx: %s\t Status: Pending   \r", tx)
				case 1:
					fmt.Printf("Tx: %s\t Status: Pending.  \r", tx)
				case 2:
					fmt.Printf("Tx: %s\t Status: Pending.. \r", tx)
				case 3:
					fmt.Printf("Tx: %s\t Status: Pending...\r", tx)
				default:
					count = -1
				}
				count++
			}
		}
	}()

	err := <-quit
	if err != nil {
		return false, HandleError(err, "", "utils.WaitForTx")
	}

	fmt.Printf("\nTx: %s\t Status: Successful\n", tx)
	return true, nil
}

// CheckBalance - check SYMBOL balance of account
func CheckBalance(address, symbol string) (float64, error) {
	url := fmt.Sprintf("http://localhost:3001/api/account/%s/balance/%s", address, symbol)

	res, err := SendRequest("GET", url, nil)
	if err != nil {
		return 0, HandleError(err, "", "utils.CheckBalance")
	}

	api, err := ControlDaemonHandler([]byte(res))
	if err != nil {
		return 0, HandleError(err, "", "utils.CheckBalance")
	}

	response := api.Response.(map[string]interface{})
	balance := response["value"].(float64)

	return balance, nil // value of $SYMBOL in account
}

// ControlDaemonHandler - handler for the API responses
func ControlDaemonHandler(_res []byte) (APIResponse, error) {
	var response = APIResponse{}

	err := json.Unmarshal(_res, &response)
	if err != nil {
		return APIResponse{}, HandleError(err, "Invalid server response", ":json.Unmarshall/utils.ControlDaemonHandler")
	}

	if !response.Success {
		return APIResponse{}, HandleError(fmt.Errorf(response.Error), response.Message, ":APIResponse/utils.ControlDaemonHandler")
	}

	return response, nil
}

// HandleError - custom error handler for the CLI.
// Uses ResponseError as a means of keeping 2 seperate error messages.
// UserMessage is a message to display to a user when an error occurs.
// LogError is a message to log or display to a developer.
// Path is the error path which is up to the developer to include.
func HandleError(err error, msg, path string) error {
	if err, ok := err.(*ErrorResponse); ok {
		return &ErrorResponse{UserMessage: err.Message() + msg, LogError: err.Error(), Path: err.Path + "/" + path}
	}
	return &ErrorResponse{UserMessage: msg, LogError: fmt.Sprint(err), Path: path}
}

// PrintError - print and logs ReponseError's.
// Use this to println the UserMessage and log the LogError with correct path.
func PrintError(err error) {
	if err, ok := err.(*ErrorResponse); ok {
		terminal.Print(ansi.Color("[ERROR] ", "196+hb"))
		terminal.Println(ansi.Color(err.Message(), "255+hb"))
		log.WithFields(log.Fields{"path": err.Path}).Fatal(err.LogError)
	} else {
		terminal.Print(ansi.Color("[ERROR] ", "196+hb"))
		terminal.Println(ansi.Color(err.Error(), "255+hb"))
		log.Fatal(err.Error())
	}
}

// GetIP - Retrieve the current machine's external IPv4 address
// using multiple ip API's.
func GetIP() (string, error) {
	sites := [4]string{"https://ipv4.myexternalip.com/raw", "https://api.ipify.org/?format=text", "https://ident.me/", "https://ipv4bot.whatismyipaddress.com"}

	for _, site := range sites {
		res, err := SendRequest("GET", site, nil)
		if err == nil {
			return res, nil
		}
	}
	return "", HandleError(fmt.Errorf("Could not retrieve IP address"), "Something went wrong getting this machines IP address", ":utils.GetIP")
}

// NewPassphrase - prompts user for new passphrase and confirms it.
func NewPassphrase() string {
	password1 := ""
	prompt := &survey.Password{
		Message: "Create a passphrase for your new wallet: ",
	}
	survey.AskOne(prompt, &password1, nil)

	password2 := ""
	prompt = &survey.Password{
		Message: "Confirm your passphrase: ",
	}
	survey.AskOne(prompt, &password2, nil)

	if strings.Compare(password1, password2) != 0 {
		fmt.Println("Passwords do not match. Please try again")
		return NewPassphrase()
	}

	return password1
}

// AskPassphrase - prompt user for passphrase.
func AskPassphrase() string {
	password := ""
	prompt := &survey.Password{
		Message: "Please type your wallet passphrase: ",
	}
	survey.AskOne(prompt, &password, nil)
	return password
}

// CachePassphrase - cache passphrase so user's don't have to retype it every
// time in the same command.
func CachePassphrase(passphrase string) {
	cachedPassphrase = passphrase
}

// Version - print version of each module
func Version() {

	res, err := SendRequest("GET", "localhost:8080/status", nil)
	if err != nil {
		PrintError(err)
	}

	terminal.Println(ansi.Color("CLI: ", "83+hb"), ansi.Color("0.4.0", "255+hb"))
	terminal.Println(ansi.Color("Control Daemon: ", "83+hb"), ansi.Color("0.4.0", "255+hb"))
	terminal.Println(ansi.Color("Network Daemon: ", "83+hb"), ansi.Color(res, "255+hb"))
}
