package keystore

import (
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
)

// CreatePGP - create a new pgp key and return path
func CreatePGP(data interface{}) (string, error) {
	url := "http://localhost:3001/api/keystore/pgp/create"

	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", fmt.Errorf("%v/keystore.CreatePGP", err)
	}

	_, err = utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/keystore.CreatePGP", err)
	}

	return fmt.Sprintf("PGP Key Created"), nil
}
