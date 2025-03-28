package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bakito/toolbox/pkg/fetcher"
)

const flagConfig = "config"

// fetchCmd represents the fetch command.
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			return err
		}
		return fetcher.New().Fetch(cfg, args...)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	addConfigFlag(fetchCmd)
}

func addConfigFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(flagConfig, "c", "",
		"The config file to be used. (default 1. '.toolbox.yaml' current dir, "+
			"2. '~/.config/toolbox.yaml', 3. '~/.toolbox.yaml')")
}
