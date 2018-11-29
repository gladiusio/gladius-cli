package main

import (
	"log"

	"github.com/gladiusio/gladius-cli/commands"
	"github.com/gladiusio/gladius-cli/config"
	"github.com/gladiusio/gladius-cli/utils"
)

// execute the command the user typed
func main() {
	// setup config handling
	config.SetupConfig("gladius-cli", config.CLIDefaults())

	// setup logger
	err := utils.SetupLogger()
	if err != nil {
		log.Fatal(err)
	}

	// execute the cmd args
	commands.Execute()
}
