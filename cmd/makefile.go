package cmd

import (
	"errors"
	"os"

	"github.com/bakito/toolbox/pkg/makefile"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

const (
	flagFile = "file"
)

// makefileCmd represents the makefile command
var makefileCmd = &cobra.Command{
	Use:   "makefile [tools]",
	Short: "Adds tools to a Makefile",
	Args: func(_ *cobra.Command, args []string) error {
		if _, err := os.Stat("tools.go"); err != nil {
			if len(args) == 0 {
				return errors.New("at least one tool must be provided")
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		mf, err := cmd.Flags().GetString(flagFile)
		if err != nil {
			return err
		}
		return makefile.Generate(client, cmd.OutOrStderr(), mf, args...)
	},
}

func init() {
	rootCmd.AddCommand(makefileCmd)

	makefileCmd.Flags().StringP(flagFile, "f", "", "The Makefile path to generate tools in")
}
