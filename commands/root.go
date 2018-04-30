package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "gladius-cli",
	Short: "CLI for Gladius Network",
	Long:  "Gladius CLI. This can be used to interact with various components of the Gladius Network.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world :)")
	},
}

// call this to "activate" commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
