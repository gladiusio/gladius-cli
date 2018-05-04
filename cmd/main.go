package main

import (
	"fmt"
	"os"

	"github.com/gladiusio/gladius-cli/commands"
	// "github.com/gladiusio/gladius-utils/config"
)

// execute the command the user typed
func main() {
	// Setup config handling
	// config.SetupConfig("test", config.CLIDefaults())
	if _, err := os.Stat("env.toml"); os.IsNotExist(err) {
		fmt.Println("env.toml not found. Please refer to the README.md to create this file")
	} else {
		commands.Execute()
	}
}
