package commands

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gladiusio/gladius-cli/keystore"
	"github.com/gladiusio/gladius-cli/node"
	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

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
	Use:   "edge [start|stop]",
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
		fmt.Println()
		terminal.Println(ansi.Color("Please add test ether to your new wallet from a ropsten faucet", "255+hb"))
		fmt.Println()
		terminal.Println(ansi.Color("Run", "255+hb"), ansi.Color("gladius create", "83+hb"), ansi.Color("again after you've acquired your test ether", "255+hb"))
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
			Name:   "email",
			Prompt: &survey.Input{Message: "What is your email?"},
			Validate: func(val interface{}) error {
				re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") // regex for email
				if val.(string) == "" {
					return errors.New("This is a required field")
				} else if !re.MatchString(val.(string)) {
					return errors.New("Please enter a valid email address")
				} else {
					return nil
				}
			},
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
		fmt.Println(err)
		return
	}

	// wait for the node tx to finish
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

	terminal.Println(ansi.Color("Node created!", "83+hb"))

	// set node data
	tx, err = node.SetNodeData(nodeAddress, answers)
	if err != nil {
		fmt.Println(err)
		return
	}

	// wait for data tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	terminal.Println(ansi.Color("Node data set!", "83+hb"))

	terminal.Print(ansi.Color("\nNode Address: ", "83+hb"))
	terminal.Print(ansi.Color(nodeAddress+"\n", "255+hb"))

	terminal.Println("\nUse", ansi.Color("gladius apply", "83+hb"), "to apply to a pool")
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

	// send data to the pool
	tx, err := node.ApplyToPool(nodeAddress, poolAddy)
	if err != nil {
		fmt.Println(err)
		return
	}

	// wait for the tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nApplication sent to pool!")
	terminal.Println("\nUse", ansi.Color("gladius check", "83+hb"), "to check the status of your application")
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

	// check application status
	status, _ := node.CheckPoolApplication(nodeAddress, poolAddy)
	fmt.Println("Pool: " + poolAddy + "\t Status: " + status)
	terminal.Println("\nUse", ansi.Color("gladius edge start", "83+hb"), "to start the edge node software")
}

// start or stop the edge daemon
func edge(cmd *cobra.Command, args []string) {

	var reply string

	if len(args) == 0 {
		fmt.Println("Please use gladius edge start or gladius edge stop")
		return
	}

	switch args[0] {
	case "start":
		reply, err := node.StartEdgeNode()
		if err != nil {
			fmt.Println("Error starting the edge node. Make sure it's running!")
		} else {
			fmt.Println("Edge Daemon:\t", reply)
			terminal.Println("\nUse", ansi.Color("gladius edge stop", "83+hb"), "to stop the edge node software")
		}
	case "stop":
		reply, err := node.StopEdgeNode()
		if err != nil {
			fmt.Println("Error stopping the edge node. Make sure it's running!")
		} else {
			fmt.Println("Edge Daemon:\t", reply)
			fmt.Println("\nUse", ansi.Color("gladius edge start", "83+hb"), "to start the edge node software")
		}
	// case "status":
	// 	reply, _ = node.StatusEdgeNode()
	default:
		reply = "command not recognized"
		fmt.Println("Edge Daemon:\t", reply)
		fmt.Println("\nUse", ansi.Color("gladius edge -h", "83+hb"), "for help")
	}
}

func test(cmd *cobra.Command, args []string) {
	address := "0x1234567890123456789012345678901234567890"
	// path := "/Users/name/.config/gladius/wallet/UTC-2018-04-14-12533634-DSFX-2234DAXF-3FSDFWEGWES.json"
	terminal.Println(ansi.Color("Wallet Address:", "83+hb"), ansi.Color(address, "255+hb"))
}

func init() {
	surveyCore.QuestionIcon = "[Gladius]"
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdEdge)
	rootCmd.AddCommand(cmdTest)
}
