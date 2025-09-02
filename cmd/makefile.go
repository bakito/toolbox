package cmd

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/bakito/toolbox/pkg/makefile"
)

const (
	flagFile    = "file"
	flagToolsGo = "tools-go"
)

var (
	toolsGo   string
	renovate  bool
	toolchain bool
	// makefileCmd represents the makefile command.
	makefileCmd = &cobra.Command{
		Use:   "makefile [tools]",
		Short: "Adds tools to a Makefile",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := resty.New()
			mf, err := cmd.Flags().GetString(flagFile)
			if err != nil {
				return err
			}
			return makefile.Generate(client, mf, renovate, toolchain, toolsGo, args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(makefileCmd)

	makefileCmd.Flags().StringP(flagFile, "f", "Makefile", "The Makefile path to generate tools in")
	makefileCmd.Flags().
		StringVar(&toolsGo, flagToolsGo, "tools.go", "The tools.go file to check for tools dependencies")
	makefileCmd.Flags().
		BoolVar(&renovate, "renovate", false, "If enabled, renovate config is added to the Makefile "+
			"(renovate.json file, if existing)")
	makefileCmd.Flags().
		BoolVar(&toolchain, "toolchain", false, "If enabled, the Makefile evaluates the go version from "+
			"go.mod and installs the tool with this go version")
}
