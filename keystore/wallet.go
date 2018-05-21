package keystore

import (
	"errors"
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// CreateWallet - create a new wallet with passphrase
func CreateWallet() error {
	url := "http://localhost:3001/api/keystore/wallet/create"

	// make a new passphrase for this wallet
	password := utils.NewPassphrase()
	pass := make(map[string]string)
	pass["passphrase"] = password

	utils.CachePassphrase(password)

	res, err := utils.SendRequest("POST", url, pass)
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateWallet", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateWallet", err)
	}

	response := api.Response.(map[string]interface{})
	address := response["address"].(string)
	path := response["path"].(string)

	fmt.Println()
	terminal.Println(ansi.Color("Wallet Address:", "83+hb"), ansi.Color(address, "255+hb"))
	fmt.Println("Wallet Path: ", path)

	return nil
}

// GetAccounts - get accounts at the standard config path
func GetAccounts() error {
	url := "http://localhost:3001/api/keystore/wallets"

	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("%v/keystore.GetAccounts", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return fmt.Errorf("%v/keystore.GetAccounts", err)
	}

	response := api.Response.([]interface{})

	if len(response) < 1 {
		return errors.New("No accounts found. Please create a wallet with: gladius-cli wallet create")
	}

	return nil
}

// EnsureAccount - Make sure they have a wallet
func EnsureAccount() (bool, error) {
	err := GetAccounts()
	if err != nil {
		return false, err
	}

	return true, nil
}
