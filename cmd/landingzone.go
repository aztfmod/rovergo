//
// Rover - Top level landing zone command
// * Doesn't do anything, all work is done by sub-commands
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

// landingzoneCmd represents the landingzone command
var landingzoneCmd = &cobra.Command{
	Use:         "landingzone",
	Aliases:     []string{"lz"},
	Short:       "Manage and deploy landing zones",
	Long:        `This command allows you to fetch landing zones or list what you have deployed`,
	Annotations: map[string]string{"cmd_group_annotation": landingzone.BuiltinCommand},
}

func init() {
	rootCmd.AddCommand(landingzoneCmd)
}
