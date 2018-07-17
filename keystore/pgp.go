package keystore

import (
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
	log "github.com/sirupsen/logrus"
)

// CreatePGP - create a new pgp key and return path
func CreatePGP(data interface{}) (string, error) {
	url := "http://localhost:3001/api/keystore/pgp/create"

	log.WithFields(log.Fields{"file": "pgp.go", "func": "CreatePGP"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", utils.HandleError(err, "", "pgp.CreatePGP")
	}

	log.WithFields(log.Fields{"file": "pgp.go", "func": "CreatePGP"}).Debug("Response recieved, piping through the response handler")
	_, err = utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", utils.HandleError(err, "", "pgp.CreatePGP")
	}

	return fmt.Sprintf("PGP Key Created"), nil
}
