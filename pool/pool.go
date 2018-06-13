package pool

import (
	"errors"
	"fmt"

	"github.com/gladiusio/gladius-cli/utils"
)

func GetOwnedPools() ([]string, error) {
	url := "http://localhost:3001/api/market/pools/owned"
	errorArray := []string{}

	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return errorArray, fmt.Errorf("%v/pool.GetOwnedPools", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return errorArray, fmt.Errorf("%v/pool.GetOwnedPools", err)
	}

	response := api.Response.([]interface{})
	responseString := make([]string, len(response))
	for i, v := range response {
		responseString[i] = v.(string)
	}

	if len(response) < 1 {
		return nil, errors.New("No owned Pools found. Please create a Pool with: gladius pool create")
	}

	return responseString, nil
}

func GetApplications(poolAddress, status string) ([]string, error) {
	url := "http://localhost:3001/api/pool/" + poolAddress + "/nodes/" + status
	errorArray := []string{}

	res, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		return errorArray, fmt.Errorf("%v/pool.GetApplications", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return errorArray, fmt.Errorf("%v/pool.GetApplications", err)
	}

	response := api.Response.([]interface{})
	addresses := make([]string, len(response))
	for i, application := range response {
		_app := application.(map[string]interface{})
		addresses[i] = _app["address"].(string)
	}

	return addresses, nil
}

func SetStatus(poolAddress, nodeAddress, status string) (string, error) {

	url := "http://localhost:3001/api/pool/" + poolAddress + "/node/" + nodeAddress + "/" + status

	res, err := utils.SendRequest("PUT", url, nil)
	if err != nil {
		return "", fmt.Errorf("%v/pool.SetStatus", err)
	}

	api, err := utils.ControlDaemonHandler([]byte(res))
	if err != nil {
		return "", fmt.Errorf("%v/pool.SetStatus", err)
	}

	response := api.Response.(map[string]interface{})
	txHash := response["txHash"].(map[string]interface{})

	return txHash["value"].(string), nil //tx hash
}
