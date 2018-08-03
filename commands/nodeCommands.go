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
		terminal.Println(ansi.Color("Your application has been sent! Use", "255+hb"), ansi.Color("gladius check", "83+hb"),
			ansi.Color("to check on the status of your application!", "255+hb"))
	}
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "applyToPool"}).Info("Application sent!")
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

	fmt.Println("Pool: " + poolAddy.(string) + "\t Status: " + status)
	terminal.Println("\nUse", ansi.Color("gladius node start", "83+hb"), "to start the node networking software if your application has been approved")
}

// start or stop the node daemon
func network(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

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
	defer utils.LogFile.Close()

	account, err := keystore.GetAccounts()
	if err != nil {
		utils.PrintError(err)
	}

	fmt.Println()
	terminal.Println(ansi.Color("Account Address:", "83+hb"), ansi.Color(account, "255+hb"))
}

func test(cmd *cobra.Command, args []string) {
	utils.SetLogLevel(utils.LogLevel)
	defer utils.LogFile.Close()

	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "Test"}).Debug("DEBUG")
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "Test"}).Info("INFO")
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "Test"}).Warning("WARNING")
	log.WithFields(log.Fields{"file": "nodeCommands.go", "func": "Test"}).Fatal("FATAL")
}

func init() {
	surveyCore.QuestionIcon = "[Gladius]"

	// register all commands
	// rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdApply)
	rootCmd.AddCommand(cmdCheck)
	rootCmd.AddCommand(cmdNetwork)
	rootCmd.AddCommand(cmdProfile)
	rootCmd.AddCommand(cmdTest)

	// register all flags
	// cmdCreate.Flags().BoolVarP(&reset, "reset", "r", false, "reset wallet")
	// rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().IntVarP(&utils.LogLevel, "level", "l", 2, "set the logging level")
	rootCmd.PersistentFlags().IntVarP(&utils.RequestTimeout, "timeout", "t", 10, "set the timeout for requests in seconds")
}
