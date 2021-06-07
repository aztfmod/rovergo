//
// Rover - Landing zone run command is the core of rover
// * This carries out actions like plan, apply or destroy via terrafrom
// * TODO: Work in progress
// * Ben C, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

var lzRunCmd = &cobra.Command{
	Use:   landingzone.ActionRun.Name(),
	Short: landingzone.ActionRun.Description(),
	Run: func(cmd *cobra.Command, args []string) {
		// Build config from command flags
		opt := landingzone.NewOptionsFromCmd(cmd)
		// And execute the relevant action
		opt.Execute(landingzone.ActionRun)
	},
}

func init() {
	// Place this command under the main `rover launchpad` command
	landingzoneCmd.AddCommand(lzRunCmd)
	landingzone.SetSharedFlags(lzRunCmd)
}
