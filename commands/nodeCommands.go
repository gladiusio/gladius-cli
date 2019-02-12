package commands

import (
	"errors"
	"fmt"
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

var cmdStatus = &cobra.Command{
	Use:   "status",
	Short: "See the status of your node",
	Long:  "See the status of each module",
	Run:   status,
}

var cmdProfile = &cobra.Command{
	Use:   "profile",
	Short: "See your profile information",
	Long:  "Display current users profile information",
	Run:   profile,
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "See the version of the Gladius Network",
	Long:  "See versions of the Gladius Network modules",
	Run:   version,
}

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start the gladius modules",
	Long:  "Start the EdgeD and Network Gateway",
	Run:   start,
}

var cmdStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop the gladius modules",
	Long:  "Stop the EdgeD and Network Gateway",
	Run:   stop,
}

var cmdUnlock = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock your wallet",
	Long:  "Unlock the gladius wallet in the Network Gateway",
	Run:   unlock,
}

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Check for updates for your node",
	Long:  "Check for updates your node modules",
	Run:   update,
}

// collect user info, send application to the server
func applyToPool(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

	// make sure they have a account, if they dont, make one
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Checking for account")
	account, _ := keystore.EnsureAccount()
	if !account {
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Warning("No account found")
		res, err := keystore.CreateAccount()
		if err != nil {
			utils.PrintError(err)
		}
		log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "createNewNode"}).Info(res)
		fmt.Println()
		terminal.Println(ansi.Color("Remember your passphrase! It's how you unlock your wallet!", "83+hb"))
		fmt.Println()
	}
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Account found")

	// create the user questions
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
		{
			Name:      "location",
			Prompt:    &survey.Input{Message: "What country are you in?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:   "estimatedSpeed",
			Prompt: &survey.Input{Message: "How much bandwidth do you have? (Mbps)"},
			Validate: func(val interface{}) error {
				re := regexp.MustCompile("^[0-9]*$") // regex for speed
				if val.(string) == "" {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Warning("Empty value")
					return errors.New("This is a required field")
				} else if !re.MatchString(val.(string)) {
					log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Warning("Invalid bandwidth value")
					return errors.New("Please enter a valid integer")
				} else {
					return nil
				}
			},
			Transform: survey.Title,
		},
		{
			Name:     "bio",
			Prompt:   &survey.Input{Message: "Why do you want to join this pool?"},
			Validate: survey.Required,
		},
	}

	// the answers will be written to this struct
	answers := make(map[string]interface{})

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Collecting application info")
	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		utils.PrintError(err)
	}

	// apply to the application server
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Sending application to server")
	_, err = node.ApplyToPool(answers["pool"].(string), answers)
	if err != nil {
		utils.PrintError(err)
	} else {
		println()
		terminal.Println(ansi.Color("Your application has been sent! Use", "255+hb"), ansi.Color("gladius check", "83+hb"),
			ansi.Color("to check on the status of your application!", "255+hb"))
	}
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Application sent!")

	checkUpdate()
}

// unlock your wallet manually
func unlock(cmd *cobra.Command, args []string) {
	utils.OpenAccount()

	checkUpdate()
}

// check the application of the node
func checkPoolApp(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

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

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Info("Collecting pool address")
	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		utils.PrintError(err)
	}

	poolAddy := answers["pool"]

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Info("Checking application")
	// check application status
	status, err := node.CheckPoolApplication(poolAddy.(string))
	if err != nil {
		utils.PrintError(err)
	}
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "checkPoolApp"}).Info("Application checked")

	fmt.Println()
	terminal.Println(ansi.Color("Pool: "+poolAddy.(string)+"\t Status: "+status, "255+hb"))
	terminal.Println(ansi.Color("\nOnce your application is approved you will automatically become an edge node!", "255+hb"))

	checkUpdate()
}

// get a users profile
func profile(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

	account, err := keystore.GetAccounts()
	if err != nil {
		utils.PrintError(err)
	}

	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(account, "255+hb"))

	checkUpdate()
}

// versions of the modules
func version(cmd *cobra.Command, args []string) {
	cli := "0.8.0"
	offline := "NOT ONLINE"

	guardian, err := node.GetVersion("guardian")
	if err != nil {
		guardian = offline
	}
	edged, err := node.GetVersion("edged")
	if err != nil {
		edged = offline
	}
	networkGateway, err := node.GetVersion("network-gateway")
	if err != nil {
		networkGateway = offline
	}

	terminal.Println(ansi.Color("CLI:", "83+hb"), ansi.Color(cli, "255+hb"))
	terminal.Println(ansi.Color("EDGED:", "83+hb"), ansi.Color(edged, "255+hb"))
	terminal.Println(ansi.Color("NETWORKD:", "83+hb"), ansi.Color(networkGateway, "255+hb"))
	terminal.Println(ansi.Color("GUARDIAN:", "83+hb"), ansi.Color(guardian, "255+hb"))

	checkUpdate()
}

func start(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

	status, err := node.Start()
	if err != nil {
		utils.PrintError(err)
	} else {
		terminal.Println(ansi.Color("Network Gateway:", "83+hb"), ansi.Color(status, "255+hb"))
		terminal.Println(ansi.Color("Edge Daemon:", "83+hb"), ansi.Color(status, "255+hb"))
	}

	checkUpdate()
}

func stop(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

	status, err := node.Stop()
	if err != nil {
		utils.PrintError(err)
	} else {
		terminal.Println(ansi.Color("Network Gateway:", "83+hb"), ansi.Color(status, "255+hb"))
		terminal.Println(ansi.Color("Edge Daemon:", "83+hb"), ansi.Color(status, "255+hb"))
	}

	checkUpdate()
}

func status(cmd *cobra.Command, args []string) {
	offline := "NOT ONLINE"
	online := "ONLINE"

	onlineColor := "83+hb"
	offlineColor := "196+hb"

	statusColor := make(map[string]string)

	_, err := node.GetVersion("guardian")
	guardian := online
	statusColor["guardian"] = onlineColor
	if err != nil {
		guardian = offline
		statusColor["guardian"] = offlineColor
	}

	_, err = node.GetVersion("edged")
	edged := online
	statusColor["edged"] = onlineColor
	if err != nil {
		edged = offline
		statusColor["edged"] = offlineColor
	}
	_, err = node.GetVersion("network-gateway")
	networkGateway := online
	statusColor["networkGateway"] = onlineColor
	if err != nil {
		networkGateway = offline
		statusColor["networkGateway"] = offlineColor
	}

	terminal.Println(ansi.Color("EDGE DAEMON:\t", statusColor["edged"]), ansi.Color(edged, "255+hb"))
	terminal.Println(ansi.Color("NETWORK GATEWAY:", statusColor["networkGateway"]), ansi.Color(networkGateway, "255+hb"))
	terminal.Println(ansi.Color("GUARDIAN:\t", statusColor["guardian"]), ansi.Color(guardian, "255+hb"))

	checkUpdate()
}

func update(cmd *cobra.Command, args []string) {
	checkUpdate()
}

func checkUpdate() {
	updateNeeded, _ := node.NeedUpdate()
	if updateNeeded {
		fmt.Println()
		fmt.Println("One or more of your modules is out of date!")
		fmt.Println("You can find the newest versions here: https://github.com/gladiusio/gladius-node")
	} else {
		fmt.Println()
		fmt.Println("Everything up to date!")
	}
}

func init() {
	surveyCore.QuestionIcon = "[Gladius]"

	// register all commands
	// rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdStatus)
	rootCmd.AddCommand(cmdProfile)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdStart)
	rootCmd.AddCommand(cmdStop)
	rootCmd.AddCommand(cmdUnlock)
	rootCmd.AddCommand(cmdUpdate)

	// register all flags
	// cmdCreate.Flags().BoolVarP(&reset, "reset", "r", false, "reset wallet")
	// rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().IntVarP(&utils.LogLevel, "level", "l", 2, "set the logging level")
	rootCmd.PersistentFlags().IntVarP(&utils.RequestTimeout, "timeout", "t", 10, "set the timeout for requests in seconds")
}
