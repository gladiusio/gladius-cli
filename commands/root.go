package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "gladius",
	Short: "CLI for Gladius Network",
	Long:  "Gladius CLI. This can be used to interact with various components of the Gladius Network.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nWelcome to the Gladius CLI!")
		fmt.Println("\nHere are the commands to create a node and apply to a pool in order:")
		fmt.Println("\n$ gladius create")
		fmt.Println("$ gladius check")
		fmt.Println("\nAfter you are accepted into a pool, you can become an edge node:")
		fmt.Println("\n$ gladius node start")
		fmt.Println("\nUse the -h flag to see the help menu")

	},
}

// Execute - call this to "activate" commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
