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

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Continuous deployment of landingzones",
	Long:  "Continuous deployment of landingzones using symphony CI/CD defintion",
}

func init() {
	rootCmd.AddCommand(cdCmd)
	cdCmd.PersistentFlags().StringP("symphony-config", "c", "", "Path to symphony config file")
	cdCmd.PersistentFlags().StringP("level", "l", "", "Level to operate on, if omitted all levels will be processed")
	_ = cobra.MarkFlagRequired(cdCmd.PersistentFlags(), "symphony-config")

	// Dynamically build sub-commands from list of actions
	for _, actionName := range landingzone.ActionEnum {
		// These actions are for CI only, makes no sense overlaping with `rover ci` here
		if actionName == landingzone.ActionFormat.Name() ||
			actionName == landingzone.ActionValidate.Name() {
			continue
		}

		action, err := landingzone.NewAction(actionName)
		cobra.CheckErr(err)
		actionSubCmd := &cobra.Command{
			Use:   action.Name(),
			Short: action.Description(),
			// Run function
			Run: func(cmd *cobra.Command, args []string) {
				levelFlag, _ := cmd.Flags().GetString("level")
				configPath, _ := cmd.Flags().GetString("symphony-config")

				// IMPORTANT: Get the action to carry out from the command name
				actionStr := cmd.Name()
				action, err := landingzone.NewAction(actionStr)
				cobra.CheckErr(err)

				conf, err := symphony.NewSymphonyConfig(configPath)
				cobra.CheckErr(err)

				if levelFlag == "" {
					conf.RunAll(action)
					return
				}

				var level *symphony.Level
				// Try to locate level in config matching the level flag passed in
				for _, confLevel := range conf.Content.Levels {
					if confLevel.Name == levelFlag {
						level = &confLevel
						break
					}
				}

				//nolint
				if level == nil {
					cobra.CheckErr(fmt.Sprintf("level '%s' not found in symphony config file", levelFlag))
				}
				//nolint
				conf.RunLevel(*level, action)
			},
		}

		// Stuff it under the parent launchpad command
		cdCmd.AddCommand(actionSubCmd)
	}
}
