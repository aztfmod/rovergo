//
// Rover - Top level launchpad command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"github.com/spf13/cobra"
)

var launchpadCmd = &cobra.Command{
	Use:     "launchpad",
	Aliases: []string{"lp"},
	Short:   "Manage and deploy launchpad, i.e. landing zone level0.",
	Long:    `Manage and deploy launchpad, i.e. landing zone level0.`,
}

func init() {
	rootCmd.AddCommand(launchpadCmd)
}
