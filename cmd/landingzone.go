//
// Rover - Top level landing zone command
// * Doesn't do anything, all work is done by sub-commands
// * Ben C, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

// landingzoneCmd represents the landingzone command
var landingzoneCmd = &cobra.Command{
	Use:     "landingzone",
	Aliases: []string{"lz"},
	Short:   "Manage and deploy landing zones",
	Long:    `This command allows you to deploy, update and destroy CAF landing zones`,
}

func init() {
	rootCmd.AddCommand(landingzoneCmd)

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
		landingzoneCmd.AddCommand(actionSubCmd)
	}
}
