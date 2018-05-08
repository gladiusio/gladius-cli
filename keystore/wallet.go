package keystore

import (
	"fmt"
	"strings"

	"github.com/gladiusio/gladius-cli/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// CreateWallet - create a new wallet with passphrase
func CreateWallet() error {
	url := "http://localhost:3001/api/keystore/wallet/create"

	password := NewPassword()

	res, err := utils.SendRequest("POST", url, fmt.Sprintf(`{"passphrase":"%s"}`, password))
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateWallet", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return fmt.Errorf("%v/keystore.CreateWallet", err)
	}

	response := api.Response.(map[string]interface{})
	address := response["address"].(string)
	path := response["address"].(string)

	fmt.Printf("Wallet Address: %s\nWallet Path: %s", address, path)

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
		println("No accounts found. Please create a wallet with: gladius-cli wallet create")
	} else {
		println("Accounts: ")
	}

	for index, element := range response {
		fmt.Printf("[%d] %s\n", index, element.(map[string]interface{})["address"].(string))
	}

	return nil
}

// NewPassword - make a new password and confirm
func NewPassword() string {
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
		return NewPassword()
	}

	return password1
}
