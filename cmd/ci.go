//
// Rover - Top level ci (continuous integration) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type SymphonyConfig struct {
	Environment  string `yaml:"environment,omitempty"`
	Repositories []struct {
		Name   string `yaml:"name,omitempty"`
		URI    string `yaml:"uri,omitempty"`
		Branch string `yaml:"branch,omitempty"`
	}
	Levels []struct {
		Level     string `yaml:"level,omitempty"`
		Type      string `yaml:"type,omitempty"`
		Launchpad bool   `yaml:"launchpad,omitempty"`
		Stacks    []struct {
			Stack             string `yaml:"stack,omitempty"`
			LandingZonePath   string `yaml:"landingZonePath,omitempty"`
			ConfigurationPath string `yaml:"configurationPath,omitempty"`
			TfState           string `yaml:"tfState,omitempty"`
		}
	}
}

var symphonyConfig SymphonyConfig

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Manage CI operations.",
	Long:  `Manage CI operations.`,
	Run: func(cmd *cobra.Command, args []string) {

		symphonyConfigFileName, _ := cmd.Flags().GetString("symphony_config")
		verbose, _ := cmd.Flags().GetBool("verbose")

		err := readAndUnmarshallConfig(symphonyConfigFileName)
		cobra.CheckErr(err)

		if verbose {
			outputVerbose(symphonyConfigFileName)
		}

		run(symphonyConfigFileName)
	},
}

func run(symphonyConfigFileName string) {
	fmt.Println()

	color.Red("Running CI command, config: %s", symphonyConfigFileName)
}

func readAndUnmarshallConfig(symphonyConfigFileName string) error {
	reader, _ := os.Open(symphonyConfigFileName)
	buf, _ := ioutil.ReadAll(reader)
	err := yaml.Unmarshal(buf, &symphonyConfig)

	return err
}

func outputVerbose(symphonyConfigFileName string) {
	fmt.Println()

	color.Blue("Verbose output of %s", symphonyConfigFileName)
	color.Green(" - Environment: %s", symphonyConfig.Environment)
	color.Green(" - Number of repositories: %d", len(symphonyConfig.Repositories))
	color.Green(" - Number of levels: %d", len(symphonyConfig.Levels))
}

func init() {
	ciCmd.Flags().StringP("symphony_config", "c", "", "Path/filename of symphony.yml")
	ciCmd.Flags().SetNormalizeFunc(aliasNormalizeFunc)

	ciCmd.Flags().BoolP("verbose", "v", false, "Output symphony.yml to console")

	err := cobra.MarkFlagRequired(ciCmd.Flags(), "symphony_config")
	cobra.CheckErr(err)

	rootCmd.AddCommand(ciCmd)
}

func aliasNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "sc":
		name = "symphony_config"
	}
	return pflag.NormalizedName(name)
}
