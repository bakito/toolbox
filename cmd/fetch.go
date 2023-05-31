package cmd

import (
	"github.com/bakito/toolbox/pkg/fetcher"
	"github.com/spf13/cobra"
)

const flagConfig = "config"

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			return err
		}
		return fetcher.New().Fetch(cfg)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringP(flagConfig, "c", "",
		"The config file to be used. (default 1. '.toolbox.yaml' current dir, 2. '~/.toolbox.yaml')")
}
