package cmd

import (
	"errors"
	"os"

	"github.com/bakito/toolbox/pkg/makefile"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

const (
	flagFile    = "file"
	flagToolsGo = "tools-go"
)

var (
	toolsGo  string
	renovate bool
	// makefileCmd represents the makefile command
	makefileCmd = &cobra.Command{
		Use:   "makefile [tools]",
		Short: "Adds tools to a Makefile",
		Args: func(_ *cobra.Command, args []string) error {
			if _, err := os.Stat(toolsGo); err != nil {
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
			return makefile.Generate(client, cmd.OutOrStderr(), mf, renovate, toolsGo, args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(makefileCmd)

	makefileCmd.Flags().StringP(flagFile, "f", "", "The Makefile path to generate tools in")
	makefileCmd.Flags().StringVar(&toolsGo, flagToolsGo, "tools.go", "The tools.go file to check for tools dependencies")
	makefileCmd.Flags().BoolVar(&renovate, "renovate", false, "If enables, renovate config is added to the Makefile (renovate.json file, if existing)")
}
