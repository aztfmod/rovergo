//
// Rover - Top level ci (continuous integration) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/symphony"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Manage CI operations.",
	Long:  `Manage CI operations.`,
	Run: func(cmd *cobra.Command, args []string) {

		symphonyConfigFileName, _ := cmd.Flags().GetString("symphony-config")
		debug, _ := cmd.Flags().GetBool("debug")

		symphonyConfig, err := symphony.NewSymphonyConfig(symphonyConfigFileName)
		cobra.CheckErr(err)

		if debug {
			symphonyConfig.OutputDebug(symphonyConfigFileName)
		}

		run(symphonyConfigFileName)
	},
}

func run(symphonyConfigFileName string) {
	fmt.Println()

	console.Infof("Running CI command, config: %s\n", symphonyConfigFileName)
}

func init() {
	ciCmd.Flags().StringP("symphony-config", "c", "", "Path/filename of symphony.yml")
	ciCmd.Flags().SetNormalizeFunc(aliasNormalizeFunc)

	ciCmd.Flags().BoolP("verbose", "v", false, "Output symphony.yml to console")

	err := cobra.MarkFlagRequired(ciCmd.Flags(), "symphony-config")
	cobra.CheckErr(err)

	rootCmd.AddCommand(ciCmd)
}

func aliasNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "sc":
		name = "symphony-config"
	}
	return pflag.NormalizedName(name)
}
