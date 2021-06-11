//
// Rover - Top level landing zone command
// * Doesn't do anything, all work is done by sub-commands
// * Ben C, May 2021
//

package cmd

import (
	"github.com/spf13/cobra"
)

// landingzoneCmd represents the landingzone command
var landingzoneCmd = &cobra.Command{
	Use:     "landingzone",
	Aliases: []string{"lz"},
	Short:   "Manage and deploy landing zones",
	Long:    `This command allows you to fetch landing zones or list what you have deployed`,
}

func init() {
	rootCmd.AddCommand(landingzoneCmd)
}
