package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/gladiusio/gladius-cli/keystore"
	"github.com/gladiusio/gladius-cli/node"
	"github.com/gladiusio/gladius-cli/utils"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var cmdEcho = &cobra.Command{
	Use:   "echo [string to echo]",
	Short: "Echo anything to the screen",
	Long: `echo is for echoing anything back.
    Echo echo’s.
    `,
	Run: echoRun,
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Deploy a new Node smart contract",
	Long:  "Deploys a new Node smart contract to the network with data",
	Run:   createNewNode,
}

var cmdApply = &cobra.Command{
	Use:   "apply",
	Short: "Apply to a Gladius Pool",
	Long:  "Send your Node's data (encrypted) to the pool owner as an application",
	Run:   applyToPool,
}

var cmdCheck = &cobra.Command{
	Use:   "check",
	Short: "Check status of your submitted pool application",
	Long:  "Check status of your submitted pool application",
	Run:   checkPoolApp,
}

var cmdEdge = &cobra.Command{
	Use:   "edge [start|stop|status]",
	Short: "Start the edge daemon",
	Long:  "Start the edge daemon networking server",
	Run:   edge,
}

var cmdTest = &cobra.Command{
	Use:   "test",
	Short: "Test function",
	Long:  "Have something to test but dont want to ruin everything else? Put it in this command!",
	Run:   test,
}

// collect user info, create node, set node data
func createNewNode(cmd *cobra.Command, args []string) {
	// make sure they have a wallet, if they dont, make one
	wallet, err := keystore.EnsureAccount()
	if !wallet {
		err = keystore.CreateWallet()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Please add test ether to your new wallet from a ropsten faucet")
		return
	}

	// create the user questions
	var qs = []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "What is your name?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:     "email",
			Prompt:   &survey.Input{Message: "What is your email?"},
			Validate: survey.Required,
		},
	}

	// the answers will be written to this struct
	answers := make(map[string]interface{})

	// perform the questions
	err = survey.Ask(qs, &answers)
	if err != nil {
		return
	}

	// gen a new pgp key for this contract
	keystore.CreatePGP(answers)

	// get ip of current machine
	ip, err := utils.GetIP()
	if err != nil {
		return
	}
	answers["ip"] = ip
	answers["status"] = "active"

	// create the node
	tx, err := node.CreateNode()
	if err != nil {
		fmt.Println("CREATE NODE: " + err.Error())
		return
	}

	// wait for the tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Node created!")

	// set node data
	tx, err = node.SetNodeData(nodeAddress, answers)
	if err != nil {
		fmt.Println(err)
		return
	}

	// wait for tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Node data set!")

	fmt.Println("\nNode Address: " + nodeAddress)
}

// send data to pool
func applyToPool(cmd *cobra.Command, args []string) {
	// make sure they have a wallet, if they dont, make one
	wallet, err := keystore.EnsureAccount()
	if !wallet {
		err = keystore.CreateWallet()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Please add test ether to your new wallet from a ropsten faucet")
		return
	}

	// build question
	poolAddy := ""
	prompt := &survey.Input{
		Message: "Pool Address: ",
	}
	survey.AskOne(prompt, &poolAddy, nil)

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	tx, err := node.ApplyToPool(nodeAddress, poolAddy)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.WaitForTx(tx)
	fmt.Println("Application sent to pool!")
}

// check the application of the node
func checkPoolApp(cmd *cobra.Command, args []string) {
	// build the prompt
	poolAddy := ""
	prompt := &survey.Input{
		Message: "Pool Address: ",
	}
	survey.AskOne(prompt, &poolAddy, nil)

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	status, _ := node.CheckPoolApplication(nodeAddress, poolAddy)
	fmt.Println("Pool: " + poolAddy + "\t Status: " + status)
}

// start - stop - status of the edge daemon
func edge(cmd *cobra.Command, args []string) {

	var reply string

	switch args[0] {
	case "start":
		reply, _ = node.StartEdgeNode()
	case "stop":
		reply, _ = node.StopEdgeNode()
	case "status":
		reply, _ = node.StatusEdgeNode()
	default:
		reply = "command not recognized"
	}
	fmt.Println("Edge Daemon:\t", reply)
}

func echoRun(cmd *cobra.Command, args []string) {
	fmt.Println(strings.Join(args, " "))
}

func test(cmd *cobra.Command, args []string) {

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan bool)

	go func() {
		i := 0
		for {
			select {
			case <-ticker.C:
				fmt.Println("REQUEST")
				if i == 4 {
					quit <- true
				}
			}
			i++
		}
	}()

	fmt.Println("Done: ", <-quit)
}

func init() {
	rootCmd.AddCommand(cmdEcho)
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdEdge)
	rootCmd.AddCommand(cmdTest)
}
