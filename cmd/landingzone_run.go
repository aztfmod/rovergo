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
	Use:   "run",
	Short: "Run actions to deploy, update or remove landing zones",

	Run: landingzone.RunCmd,
}

func init() {
	landingzoneCmd.AddCommand(lzRunCmd)
	landingzone.SetSharedFlags(lzRunCmd)
}
