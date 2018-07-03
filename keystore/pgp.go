package keystore

import (
	"fmt"
	"os"

	"github.com/gladiusio/gladius-cli/utils"
	log "github.com/sirupsen/logrus"
)

// CreatePGP - create a new pgp key and return path
func CreatePGP(data interface{}) (string, error) {
	url := "http://localhost:3001/api/keystore/pgp/create"

	log.WithFields(log.Fields{"file": "pgp.go", "func": "CreatePGP"}).Debug("POST: ", url)
	res, err := utils.SendRequest("POST", url, data)
	if err != nil {
		return "", fmt.Errorf("%v/keystore.CreatePGP", err)
	}

	log.WithFields(log.Fields{"file": "pgp.go", "func": "CreatePGP"}).Debug("Response recieved, piping through the response handler")
	_, err = utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/keystore.CreatePGP", err)
	}

	return fmt.Sprintf("PGP Key Created"), nil
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
