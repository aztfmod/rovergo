//
// Rover - Top level launchpad command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

var launchpadCmd = &cobra.Command{
	Use:     "launchpad",
	Aliases: []string{"lp"},
	Short:   "Manage and deploy a launchpad, i.e. landing zone level0.",
	Long:    `Manage and deploy a launchpad, i.e. landing zone level0.`,
}

func init() {
	rootCmd.AddCommand(launchpadCmd)

	// Dynamically build sub-commands from list of actions
	for _, actionName := range landingzone.ActionEnum {
		action, err := landingzone.NewAction(actionName)
		cobra.CheckErr(err)
		actionSubCmd := &cobra.Command{
			Use:   action.Name(),
			Short: action.Description(),
			Run: func(cmd *cobra.Command, args []string) {
				// Build config from command flags
				opt := landingzone.NewOptionsFromCmd(cmd)
				// And execute the relevant action
				opt.Execute(action)

				console.Success("Rover has finished")
			},
		}
		// Set all the shared action flags
		landingzone.SetSharedFlags(actionSubCmd)
		// Stuff it under the parent launchpad command
		launchpadCmd.AddCommand(actionSubCmd)
	}
}
