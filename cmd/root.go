// Package cmd implements cobra commands
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bakito/toolbox/version"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     "toolbox",
	Short:   "🧰 a small toolbox helping to fetch tools",
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
