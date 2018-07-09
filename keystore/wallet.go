package keystore

import (
	"errors"
	"fmt"
	"os"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	log "github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// LogFile - Where the logs are stored
var LogFile *os.File

// CreateAccount - create a new account with passphrase
func CreateAccount() error {
	url := "http://localhost:3001/api/keystore/account/create"

	// make a new passphrase for this account
	password := utils.NewPassphrase()
	pass := make(map[string]string)
	pass["passphrase"] = password

	utils.CachePassphrase(password)
	log.WithFields(log.Fields{"file": "wallet.go", "func": "CreateAccount"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, pass)
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateAccount", err)
	}

	println(res)

	log.WithFields(log.Fields{"file": "wallet.go", "func": "CreateAccount"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateAccount", err)
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "CreateAccount"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	address := response["address"].(string)

	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(address, "255+hb"))

	return nil
}

// GetAccounts - get accounts at the standard config path
func GetAccounts() (string, error) {
	url := "http://localhost:3001/api/keystore/account"

	log.WithFields(log.Fields{"file": "wallet.go", "func": "GetAccounts"}).Debug("POST: ", url)
	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/keystore.GetAccounts", err)
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "GetAccounts"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/keystore.GetAccounts", err)
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "GetAccounts"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})

	if len(response) < 1 {
		return "", errors.New("No accounts found/keystore.GetAccounts")
	}

	return response["address"].(string), nil
}

// EnsureAccount - Make sure they have an account
func EnsureAccount() (bool, error) {
	_, err := GetAccounts()
	if err != nil {
		return false, err
	}

	return true, nil
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
