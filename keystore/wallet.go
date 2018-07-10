package keystore

import (
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	log "github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// CreateAccount - create a new account with passphrase
func CreateAccount() (string, error) {
	url := "http://localhost:3001/api/keystore/account/create"

	// make a new passphrase for this account
	password := utils.NewPassphrase()
	pass := make(map[string]string)
	pass["passphrase"] = password

	utils.CachePassphrase(password)
	log.WithFields(log.Fields{"file": "wallet.go", "func": "CreateAccount"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, pass)
	if err != nil {
		return "", utils.HandleError(err, "", "wallet.CreateAccount")
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "CreateAccount"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", utils.HandleError(err, "", "walle.CreateAccount")
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "CreateAccount"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	address := response["address"].(string)

	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(address, "255+hb"))

	return "Account created", nil
}

// GetAccounts - get accounts at the standard config path
func GetAccounts() (string, error) {
	url := "http://localhost:3001/api/keystore/account"

	log.WithFields(log.Fields{"file": "wallet.go", "func": "GetAccounts"}).Debug("POST: ", url)
	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return "", utils.HandleError(err, "", "wallet.GetAccounts")
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "GetAccounts"}).Debug("Response recieved, piping through the response handler")
	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", utils.HandleError(err, "", "wallet.GetAccounts")
	}

	log.WithFields(log.Fields{"file": "wallet.go", "func": "GetAccounts"}).Debug("Decoding response fields")
	response := api.Response.(map[string]interface{})
	if len(response) < 1 {
		return "", utils.HandleError(err, "", "wallet.GetAccounts")
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
