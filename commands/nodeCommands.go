package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gladiusio/gladius-cli/internal"
	"github.com/gladiusio/gladius-cli/node"
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
	answers := node.Node{}

	// perform the questions
	err := survey.Ask(qs, &answers.Data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// need to collect ip
	answers.Data.IPAddress = "1.1.1.1"
	answers.Data.Status = "active"

	// save the struct to a file (im gonna turn these into go routines and hopefully find a good way to condense these lines)
	if err = utils.WriteToEnv("node", "type", "node", "env.toml", "env.toml"); err != nil {
		fmt.Println(err)
		return
	}
	if err = utils.WriteToEnv("node", "status", answers.Data.Status, "env.toml", "env.toml"); err != nil {
		fmt.Println(err)
		return
	}
	if err = utils.WriteToEnv("node", "name", answers.Data.Name, "env.toml", "env.toml"); err != nil {
		fmt.Println(err)
		return
	}
	if err = utils.WriteToEnv("node", "email", answers.Data.Email, "env.toml", "env.toml"); err != nil {
		fmt.Println(err)
		return
	}
	if err = utils.WriteToEnv("node", "ipAddress", answers.Data.IPAddress, "env.toml", "env.toml"); err != nil {
		fmt.Println(err)
		return
	}
	if err = utils.WriteToEnv("node", "status", answers.Data.Status, "env.toml", "env.toml"); err != nil {
		fmt.Println(err)
		return
	}

	// create the node
	tx, err := node.CreateNode()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("")
	node.WaitForTx(tx)

	// save the node address
	nodeAddress := node.GetNodeAddress()
	if err = utils.WriteToEnv("node", "address", nodeAddress, "env.toml", "env.toml"); err != nil {
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

	node.WaitForTx(tx)
	fmt.Println("Node data set!")

	fmt.Println("\n" + nodeAddress)
}

// send data to pool
func applyToPool(cmd *cobra.Command, args []string) {
	envFile, err := utils.GetEnvMap("env.toml")
	envNode := envFile["node"]
	env := envFile["environment"]

	fmt.Println(env["poolAddress"])

	// build question
	poolAddy := ""
	prompt := &survey.Input{
		Message: "Pool Address: ",
	}
	survey.AskOne(prompt, &poolAddy, nil)

	tx, err := node.ApplyToPool(envNode["address"], poolAddy)
	if err != nil {
		fmt.Println(err)
		return
	}

	node.WaitForTx(tx)
	fmt.Println("Application sent to pool!")
}

// check the application of the node
func checkPoolApp(cmd *cobra.Command, args []string) {
	envFile, err := utils.GetEnvMap("env.toml")
	if err != nil {
		fmt.Println(err)
		return
	}

	envNode := envFile["node"]
	env := envFile["environment"]

	fmt.Println(env["poolAddress"])

	// build the prompt
	poolAddy := ""
	prompt := &survey.Input{
		Message: "Pool Address: ",
	}
	survey.AskOne(prompt, &poolAddy, nil)

	status := node.CheckPoolApplication(envNode["address"], poolAddy)
	fmt.Println("Pool: " + poolAddy + "\t Status: " + status)
}

// start - stop - status of the edge daemon
func edge(cmd *cobra.Command, args []string) {

	var reply string
	switch args[0] {
	case "start":
		reply = node.StartEdgeNode()
	case "stop":
		reply = node.StopEdgeNode()
	case "status":
		reply = node.StatusEdgeNode()
	default:
		reply = "command not recognized"
	}
	fmt.Println("Edge Daemon:\t", reply)
}

func echoRun(cmd *cobra.Command, args []string) {
	fmt.Println(strings.Join(args, " "))
}

func test(cmd *cobra.Command, args []string) {
	b, err := ioutil.ReadFile("env.toml") // read env file
	if err != nil {
		fmt.Println("Error reading: " + "env.toml")
	}

	var envFile = make(map[string]map[string]string)

	if _, err := toml.Decode(string(b), &envFile); err != nil { // turn file into mapping
		fmt.Println("Error decoding")
	}

	envFile["node"]["hello"] = "test"

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(envFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	err = ioutil.WriteFile("env.toml", (*buf).Bytes(), 0644)

	// utils.WriteToEnv("node", "another", "test", "env.toml", "env.toml")

}

func init() {
	node.PostSettings("env.toml")

	rootCmd.AddCommand(cmdEcho)
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdEdge)
	rootCmd.AddCommand(cmdTest)
}
