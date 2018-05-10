package keystore

import (
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
)

func CreatePGP(data interface{}) (string, error) {
	url := "http://localhost:3001/api/keystore/pgp/create"

	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", fmt.Errorf("%v/keystore.CreatePGP", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/keystore.CreatePGP", err)
	}

	response := api.Response.(map[string]interface{})
	path := response["path"].(string)

	return fmt.Sprintf("PGP Path: %s", path), nil
}
