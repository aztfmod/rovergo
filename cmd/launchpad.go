//
// Rover - Top level launchpad command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// launchpadCmd represents the launchpad command
var launchpadCmd = &cobra.Command{
	Use:     "launchpad",
	Aliases: []string{"lp"},
	Short:   "Manage and deploy launchpad, i.e. landing zone level0.",
	Long:    `Manage and deploy launchpad, i.e. landing zone level0.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("launchpad called")
	},
}

func init() {
	rootCmd.AddCommand(launchpadCmd)
}
