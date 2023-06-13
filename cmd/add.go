package cmd

import (
	"fmt"
	"log"

	"github.com/bakito/toolbox/pkg/fetcher"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/spf13/cobra"
)

const (
	flagAdditional  = "additional"
	flagGithub      = "github"
	flagGoogle      = "google"
	flagDownloadURL = "downloadURL"
	flagVersion     = "version"
)

// fetchCmd represents the fetch command
var addCmd = &cobra.Command{
	Use:   "add <tool-name>",
	Short: "Add a tool to the config",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		github, err := cmd.Flags().GetString(flagGithub)
		if err != nil {
			return err
		}
		google, err := cmd.Flags().GetString(flagGoogle)
		if err != nil {
			return err
		}
		downloadURL, err := cmd.Flags().GetString(flagDownloadURL)
		if err != nil {
			return err
		}

		if github == "" && downloadURL == "" && google == "" {
			return fmt.Errorf("either %q, %q or %q must be defined", flagGithub, flagDownloadURL, flagGoogle)
		}

		cfg, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			return err
		}
		toolbox, filePath, err := fetcher.ReadToolbox(cfg)
		if err != nil {
			return err
		}

		version, err := cmd.Flags().GetString(flagVersion)
		if err != nil {
			return err
		}

		additional, err := cmd.Flags().GetStringArray(flagAdditional)
		if err != nil {
			return err
		}

		if tool, ok := toolbox.Tools[args[0]]; ok {
			log.Printf("🔄 Updating tool %s\n", args[0])
			if google != "" {
				tool.Github = ""
				tool.Google = google
				tool.DownloadURL = ""
			} else if downloadURL != "" {
				tool.Github = ""
				tool.Google = ""
				tool.DownloadURL = downloadURL
			} else {
				tool.Github = github
				tool.Google = ""
				tool.DownloadURL = ""
			}
			if version != "" {
				tool.Version = version
			}
			if len(additional) != 0 {
				tool.Additional = additional
			}
		} else {
			log.Printf("↘️ Adding tool %s\n", args[0])
			toolbox.Tools[args[0]] = &types.Tool{
				Github:      github,
				DownloadURL: downloadURL,
				Version:     version,
				Additional:  additional,
			}
		}
		log.Println("💾 Saving config")
		return fetcher.SaveYamlFile(filePath, toolbox)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addConfigFlag(addCmd)
	addCmd.Flags().String(flagGithub, "", "The tool's github repo")
	addCmd.Flags().String(flagGoogle, "", "The tool's google URL")
	addCmd.Flags().String(flagDownloadURL, "", "The tool's download URL")
	addCmd.Flags().String(flagVersion, "", "The tool's version or version URL")
	addCmd.Flags().StringArray(flagAdditional, nil, "Additional tools to be fetched")
}
