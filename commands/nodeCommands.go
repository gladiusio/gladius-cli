package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/gladiusio/gladius-cli/keystore"
	"github.com/gladiusio/gladius-cli/node"
	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

var cmdCreate = &cobra.Command{
	Use:   "create [wallet or node] [name(node)] [email(node)]",
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

var cmdNetwork = &cobra.Command{
	Use:   "node [start|stop|status]",
	Short: "Start/Stop or check status of your node's networking server",
	Long:  "Start/Stop or check status of your node's networking server",
	Run:   network,
}

var cmdProfile = &cobra.Command{
	Use:   "profile",
	Short: "See your profile information",
	Long:  "Display current users profile information",
	Run:   profile,
}

var cmdAccount = &cobra.Command{
	Use:   "account",
	Short: "Print users account address",
	Long:  "Print users account address",
	Run:   account,
}

var cmdTest = &cobra.Command{
	Use:   "test",
	Short: "Test function",
	Long:  "Have something to test but dont want to ruin everything else? Put it in this command!",
	Run:   test,
}

// collect user info, create node, set node data
func createNewNode(cmd *cobra.Command, args []string) {
	// user data passed in
	answers := make(map[string]interface{})

	nodeType := true

	// only accept 3 args: type, name, and email
	for index, arg := range args {
		if index >= 3 {
			break
		}
		if index == 0 {
			if strings.Compare(arg, "wallet") == 0 {
				nodeType = false
			}
		}
		switch index {
		case 1:
			answers["name"] = arg
		case 2:
			answers["email"] = arg
		}
	}

	// if they do gladius create wallet
	if !nodeType {
		err := keystore.CreateWallet()
		if err != nil {
			panic(err)
		}
		return
	}

	// gen a new pgp key for this contract
	keystore.CreatePGP(answers)

	// get ip of current machine
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	// everything before the query
	ip := strings.Split(addrs[4].String(), "/")[0]

	answers["ip"] = ip

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
	// user data passed in
	answers := make(map[string]interface{})

	// only accept 1 arg1: pool
	for index, arg := range args {
		if index >= 1 {
			break
		}
		switch index {
		case 0:
			answers["pool"] = arg
		}
	}

	poolAddy := answers["pool"]

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	// send data to the pool
	tx, err := node.ApplyToPool(nodeAddress, poolAddy.(string))
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
	// user data passed in
	answers := make(map[string]interface{})

	// only accept 1 arg1: pool
	for index, arg := range args {
		if index >= 1 {
			break
		}
		switch index {
		case 0:
			answers["pool"] = arg
		}
	}

	poolAddy := answers["pool"]

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	// check application status
	status, _ := node.CheckPoolApplication(nodeAddress, poolAddy.(string))
	fmt.Println("Pool: " + poolAddy.(string) + "\t Status: " + status)
	terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software")
}

// start or stop the node daemon
func network(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please use: \ngladius node start\ngladius node stop\ngladius node status")
		return
	}

	switch args[0] {
	case "start":
		reply, err := node.StartNetworkNode()
		if err != nil {
			fmt.Println("Error starting the node networking daemon. Make sure it's running!")
		} else {
			terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
			terminal.Println("\nUse", ansi.Color("gladius node stop", "83+hb"), "to stop the node networking software")
			terminal.Println("Use", ansi.Color("gladius node status", "83+hb"), "to check the status of the node networking software")
		}
	case "stop":
		reply, err := node.StopNetworkNode()
		if err != nil {
			fmt.Println("Error stopping the node networking daemon. Make sure it's running!")
		} else {
			terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
			terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software")
			terminal.Println("Use", ansi.Color("gladius node status", "83+hb"), "to check the status of the node networking software")
		}
	case "status":
		reply, err := node.StatusNetworkNode()
		if err != nil {
			fmt.Println("Error communicating with the node networking daemon. Make sure it's running!")
		} else {
			terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
			terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software")
			terminal.Println("Use", ansi.Color("gladius node stop", "83+hb"), "to stop the node networking software")
		}
	default:
		reply := "command not recognized"
		terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
		terminal.Println("\nUse", ansi.Color("gladius node -h", "83+hb"), "for help")
	}
}

// get a users profile
func profile(cmd *cobra.Command, args []string) {
	accounts, err := keystore.GetAccounts()
	if err != nil {
		fmt.Println("No accounts found. Create a wallet with: gladius create")
		return
	}
	wallet := accounts[0].(map[string]interface{})
	userAddress := wallet["address"].(string)
	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(userAddress, "255+hb"))

	address, err := node.GetNodeAddress()
	if err != nil {
		fmt.Println("No Node found. Create a node with : gladius create")
		return
	}
	terminal.Println(ansi.Color("Node Address:", "83+hb"), ansi.Color(address, "255+hb"))

	data, err := node.GetNodeData(address)
	if err != nil {
		fmt.Println("No Node found. Create a node with : gladius create")
		return
	}
	terminal.Println(ansi.Color("Node Name:", "83+hb"), ansi.Color(data["name"].(string), "255+hb"))
	terminal.Println(ansi.Color("Node Email:", "83+hb"), ansi.Color(data["email"].(string), "255+hb"))
	terminal.Println(ansi.Color("Node IP:", "83+hb"), ansi.Color(data["ip"].(string), "255+hb"))
}

// print users wallet address
func account(cmd *cobra.Command, args []string) {
	accounts, err := keystore.GetAccounts()
	if err != nil {
		fmt.Println("No accounts found. Create a wallet with: gladius create")
		return
	}
	wallet := accounts[0].(map[string]interface{})
	userAddress := wallet["address"].(string)
	fmt.Println(userAddress)
}

// test function
func test(cmd *cobra.Command, args []string) {
	// for _, arg := range args {
	// 	fmt.Println(arg)
	// }
	for index, arg := range args {
		if index >= 2 {
			break
		} else {
			fmt.Printf("%d %v \n", index, arg)
		}
	}
}

func init() {
	surveyCore.QuestionIcon = "[Gladius]"
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdNetwork)
	rootCmd.AddCommand(cmdProfile)
	rootCmd.AddCommand(cmdAccount)
	rootCmd.AddCommand(cmdTest)
}
