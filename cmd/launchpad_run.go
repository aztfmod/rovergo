//
// Rover - Launchpad run command
// * This carries out actions like plan, apply or destroy via terrafrom
// * TODO: Work in progress
// * Ben C, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

var lpRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run actions to deploy, update or remove launchpads",

	Run: landingzone.RunFunc,
}

func init() {
	launchpadCmd.AddCommand(lpRunCmd)
	landingzone.SetSharedFlags(lpRunCmd)
}
