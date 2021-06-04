//
// Rover - Top level cd (continuous deployment) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/symphony"
	"github.com/spf13/cobra"
)

const allLevels = -1

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Manage CD operations.",
	Long:  `Manage CD operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		levelInt, _ := cmd.Flags().GetInt("level")
		configPath, _ := cmd.Flags().GetString("symphony-config")
		actionStr, _ := cmd.Flags().GetString("action")
		action, err := landingzone.NewAction(actionStr)
		cobra.CheckErr(err)

		conf, err := symphony.NewSymphonyConfig(configPath)
		cobra.CheckErr(err)

		if levelInt == allLevels {
			conf.RunAll(action)
			return
		}

		if levelInt < allLevels {
			cobra.CheckErr("Level must be greater than zero")
		}

		var level *symphony.Level
		// Try to locate level in config from int
		for _, confLevel := range conf.Levels {
			if confLevel.Number == levelInt {
				level = &confLevel
				break
			}
		}

		//nolint
		if level == nil {
			cobra.CheckErr(fmt.Sprintf("level '%d' not found in symphony config file", levelInt))
		}
		//nolint
		conf.RunLevel(*level, action)
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
	cdCmd.Flags().StringP("symphony-config", "c", "", "Path to symphony config file")
	cdCmd.Flags().IntP("level", "l", allLevels, "Level to operate on")
	landingzone.AddActionFlag(cdCmd)

	_ = cobra.MarkFlagRequired(cdCmd.Flags(), "symphony-config")
}
