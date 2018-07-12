package commands

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/gladiusio/gladius-cli/keystore"
	"github.com/gladiusio/gladius-cli/node"
	"github.com/gladiusio/gladius-cli/utils"
	"github.com/mgutz/ansi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// LogFile - Where the logs are stored
var LogFile *os.File

// var reset bool

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

var cmdTest = &cobra.Command{
	Use:   "test",
	Short: "Test function",
	Long:  "Have something to test but dont want to ruin everything else? Put it in this command!",
	Run:   test,
}

// collect user info, create node, set node data
func createNewNode(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer LogFile.Close()

	// make sure they have a account, if they dont, make one
	account, _ := keystore.EnsureAccount()
	if !account {
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("No account found")
		res, err := keystore.CreateAccount()
		if err != nil {
			utils.PrintError(err)
		}
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info(res)
		fmt.Println()
		terminal.Println(ansi.Color("Please add test ether to your new account from a ropsten faucet", "255+hb"))
		fmt.Println()
		terminal.Println(ansi.Color("Run", "255+hb"), ansi.Color("gladius create", "83+hb"), ansi.Color("again after you've acquired your test ether", "255+hb"))
		return
	}

	acc, err := keystore.GetAccounts()
	if err != nil {
		utils.PrintError(err)
	}

	// check balance before they start
	balance, err := utils.CheckBalance(acc, "eth")
	if balance < 0.1 {
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("Account balance low (ETH)")
	}
	if balance == 0 {
		zero := utils.HandleError(fmt.Errorf("Account has no ether"), "You have no funds! Please add test ether to your account!", "nodeCommands.createNewNode")
		utils.PrintError(zero)
	}
	if err != nil {
		utils.PrintError(err)
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
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("Empty value")
					return errors.New("This is a required field")
				} else if !re.MatchString(val.(string)) {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("Invalid Email")
					return errors.New("Please enter a valid email address")
				} else {
					return nil
				}
			},
		},
	}

	// get ip address and if sites are down then make user enter manually
	ipSuccess := false
	ip, err := utils.GetIP()
	if err != nil {
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning(err)
		qs = []*survey.Question{
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
						log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("Empty value")
						return errors.New("This is a required field")
					} else if !re.MatchString(val.(string)) {
						log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("Invalid email")
						return errors.New("Please enter a valid email address")
					} else {
						return nil
					}
				},
			},
			{
				Name:      "ip",
				Prompt:    &survey.Input{Message: "What is your ip address?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
		}
	}

	// the answers will be written to this struct
	answers := make(map[string]interface{})

	// perform the questions
	err = survey.Ask(qs, &answers)
	if err != nil {
		utils.PrintError(err)
	}

	// gen a new pgp key for this contract
	res, err := keystore.CreatePGP(answers)
	if err != nil {
		utils.PrintError(err)
	}
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info(res)

	if !ipSuccess {
		answers["ip"] = ip
	}

	answers["status"] = "active"

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info("Creating Node")
	// create the node
	tx, err := node.CreateNode()
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info("Waiting for TX")
	// wait for the node tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		utils.PrintError(err)
	}

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info("Node created")
	terminal.Println(ansi.Color("Node created!", "83+hb"))
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info("Setting Node data")

	// set node data
	tx, err = node.SetNodeData(nodeAddress, answers)
	if err != nil {
		utils.PrintError(err)
	}

	// wait for data tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info("Node data set")
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info("Node fully created")

	terminal.Println(ansi.Color("Node data set!", "83+hb"))

	terminal.Print(ansi.Color("\nNode Address: ", "83+hb"))
	terminal.Print(ansi.Color(nodeAddress+"\n", "255+hb"))

	terminal.Println("\nUse", ansi.Color("gladius apply", "83+hb"), "to apply to a pool")
}

// send data to pool
func applyToPool(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer LogFile.Close()

	// make sure they have a account, if they dont, make one
	account, _ := keystore.EnsureAccount()
	if !account {
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Warning("No accounts found")
		res, err := keystore.CreateAccount()
		if err != nil {
			utils.PrintError(err)
		}
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info(res)
		fmt.Println("Please add test ether to your new account from a ropsten faucet")
		return
	}

	// build question
	var qs = []*survey.Question{
		{
			Name:   "pool",
			Prompt: &survey.Input{Message: "Pool Address: "},
			Validate: func(val interface{}) error {
				re := regexp.MustCompile("^0x[a-fA-F0-9]{40}$") // regex for email
				if val.(string) == "" {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Warning("Empty value")
					return errors.New("This is a required field")
				} else if !re.MatchString(val.(string)) {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Warning("Invalid email")
					return errors.New("Please enter a valid ethereum address")
				} else {
					return nil
				}
			},
		},
	}

	// the answers will be written to this struct
	answers := make(map[string]interface{})

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		utils.PrintError(err)
	}

	poolAddy := answers["pool"]

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Applying to pool")
	// send data to the pool
	tx, err := node.ApplyToPool(nodeAddress, poolAddy.(string))
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Waiting for TX")
	// wait for the tx to finish
	_, err = utils.WaitForTx(tx)
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Application transaction successful")
	fmt.Println("\nApplication sent to pool!")
	terminal.Println("\nUse", ansi.Color("gladius check", "83+hb"), "to check the status of your application")
}

// check the application of the node
func checkPoolApp(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer LogFile.Close()

	// build question
	var qs = []*survey.Question{
		{
			Name:   "pool",
			Prompt: &survey.Input{Message: "Pool Address: "},
			Validate: func(val interface{}) error {
				re := regexp.MustCompile("^0x[a-fA-F0-9]{40}$") // regex for email
				if val.(string) == "" {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Warning("Empty value")
					return errors.New("This is a required field")
				} else if !re.MatchString(val.(string)) {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Warning("Invalid ETH address")
					return errors.New("Please enter a valid ethereum address")
				} else {
					return nil
				}
			},
		},
	}

	// the answers will be written to this struct
	answers := make(map[string]interface{})

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		utils.PrintError(err)
	}

	poolAddy := answers["pool"]

	// save the node address
	nodeAddress, err := node.GetNodeAddress()
	if err != nil {
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Info("Checking Application")
	// check application status
	status, _ := node.CheckPoolApplication(nodeAddress, poolAddy.(string))
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Info("Application checked")
	fmt.Println("Pool: " + poolAddy.(string) + "\t Status: " + status)
	terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software")
}

// start or stop the node daemon
func network(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer LogFile.Close()

	if len(args) == 0 {
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "network"}).Fatal("Please use: \ngladius node start\ngladius node stop\ngladius node status")
	}

	switch args[0] {
	case "start":
		reply, err := node.StartNetworkNode()
		if err != nil {
			utils.PrintError(err)
		} else {
			log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "network"}).Info("Network daemon started")
			terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
			terminal.Println("\nUse", ansi.Color("gladius node stop", "83+hb"), "to stop the node networking software")
			terminal.Println("Use", ansi.Color("gladius node status", "83+hb"), "to check the status of the node networking software")
		}
	case "stop":
		reply, err := node.StopNetworkNode()
		if err != nil {
			utils.PrintError(err)
		} else {
			log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "network"}).Info("Network daemon stopped")
			terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
			terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software")
			terminal.Println("Use", ansi.Color("gladius node status", "83+hb"), "to check the status of the node networking software")
		}
	case "status":
		reply, err := node.StatusNetworkNode()
		if err != nil {
			utils.PrintError(err)
		} else {
			log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "network"}).Info("Network daemon status")
			terminal.Println(ansi.Color("Network Daemon:\t", "83+hb"), ansi.Color(reply, "255+hb"))
			terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software")
			terminal.Println("Use", ansi.Color("gladius node stop", "83+hb"), "to stop the node networking software")
		}
	default:
		terminal.Println("\nUse", ansi.Color("gladius node -h", "83+hb"), "for help")
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "network"}).Fatal("command not recognized")
	}
}

// get a users profile
func profile(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer LogFile.Close()

	account, err := keystore.GetAccounts()
	if err != nil {
		utils.PrintError(err)
		// fmt.Println("No accounts found. Create a account with: gladius create")
	}

	userAddress := account
	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(userAddress, "255+hb"))

	address, err := node.GetNodeAddress()
	if err != nil {
		// fmt.Println("No Node found. Create a node with : gladius create")
		utils.PrintError(err)
	}

	terminal.Println(ansi.Color("Node Address:", "83+hb"), ansi.Color(address, "255+hb"))

	data, err := node.GetNodeData(address)
	if err != nil {
		// fmt.Println("No Node found. Create a node with : gladius create")
		utils.PrintError(err)
	}

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "profile"}).Info("Node information found")
	terminal.Println(ansi.Color("Node Name:", "83+hb"), ansi.Color(data["name"].(string), "255+hb"))
	terminal.Println(ansi.Color("Node Email:", "83+hb"), ansi.Color(data["email"].(string), "255+hb"))
	terminal.Println(ansi.Color("Node IP:", "83+hb"), ansi.Color(data["ip"].(string), "255+hb"))
}

func test(cmd *cobra.Command, args []string) {
	customError := utils.ErrorResponse{UserMessage: "MSG", LogError: "LOG", Path: "PATH"}
	err := utils.HandleError(&customError, "msg2", "path2")
	utils.PrintError(err)
}

func init() {
	surveyCore.QuestionIcon = "[Gladius]"

	// register all commands
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdNetwork)
	rootCmd.AddCommand(cmdProfile)
	rootCmd.AddCommand(cmdTest)

	// register all flags
	// cmdCreate.Flags().BoolVarP(&reset, "reset", "r", false, "reset wallet")
	// rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().IntVarP(&utils.LogLevel, "level", "l", 2, "set the logging level")

	// clear previous log file
	utils.ClearLogger()

	LogFile, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Warning("Failed to log to file, using default stderr")
	}

	log.SetOutput(LogFile)

}
