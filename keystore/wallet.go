package keystore

import (
	"errors"
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// CreateAccount - create a new account with passphrase
func CreateAccount() error {
	url := "http://localhost:3001/api/keystore/account/create"

	// make a new passphrase for this account
	password := utils.NewPassphrase()
	pass := make(map[string]string)
	pass["passphrase"] = password

	utils.CachePassphrase(password)

	res, err := utils.SendRequest("POST", url, pass)
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateAccount", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateAccount", err)
	}

	response := api.Response.(map[string]interface{})
	address := response["address"].(string)

	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(address, "255+hb"))

	return nil
}

// GetAccounts - get accounts at the standard config path
func GetAccounts() ([]interface{}, error) {
	url := "http://localhost:3001/api/keystore/account"

	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%v/keystore.GetAccounts", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return nil, fmt.Errorf("%v/keystore.GetAccounts", err)
	}

	response := api.Response.([]interface{})

	if len(response) < 1 {
		return nil, errors.New("No accounts found/keystore.GetAccounts")
	}

	return response, nil
}

// EnsureAccount - Make sure they have an account
func EnsureAccount() (bool, error) {
	_, err := GetAccounts()
	if err != nil {
		return false, err
	}

	return true, nil
}
