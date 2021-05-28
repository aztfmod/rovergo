//
// Rover - Landing zone run command is the core of rover
// * This carries out actions like plan, apply or destroy via terrafrom
// * TODO: !!THIS IS NOT EVEN CLOSE TO BEING COMPLETED!!
// * Ben C, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/core"
	"github.com/spf13/cobra"
)

var lpRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an action for launchpad",

	Run: core.RunCmd,
}

func init() {
	launchpadCmd.AddCommand(lpRunCmd)
	core.SetSharedFlags(lpRunCmd)
}
