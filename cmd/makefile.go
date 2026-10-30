package cmd

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/bakito/toolbox/pkg/makefile"
)

const (
	flagFile      = "file"
	flagToolsGo   = "tools-go"
	flagToolchain = "toolchain"
)

var (
	toolsGo   string
	renovate  bool
	toolchain string
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
			goVersion := toolchain
			if goVersion == "go.mod" {
				goVersion = ""
			}
			return makefile.Generate(client, mf, renovate, cmd.Flags().Changed(flagToolchain), goVersion, toolsGo, args...)
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
		StringVar(&toolchain, flagToolchain, "", "If set, use the specified Go version (numeric semver, e.g. 1.27.3); "+
			"a bare flag reads the version from go.mod")
	makefileCmd.Flags().Lookup(flagToolchain).NoOptDefVal = "go.mod"
}
