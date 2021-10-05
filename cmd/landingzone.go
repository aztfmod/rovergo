//
// Rover - Top level landing zone command
// * Doesn't do anything, all work is done by sub-commands
//

package cmd

import (
	"fmt"

	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

// landingzoneCmd represents the landingzone command
var landingzoneCmd = &cobra.Command{
	Use:     "landingzone",
	Aliases: []string{"lz"},
	Short:   fmt.Sprintf("[%s command]\tManage and deploy landing zones", landingzone.BuiltinCommand),
	Long:    `This command allows you to fetch landing zones or list what you have deployed`,
}

func init() {
	rootCmd.AddCommand(landingzoneCmd)
}
