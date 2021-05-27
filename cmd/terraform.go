//
// Rover - Top level terraform command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:     "terraform",
	Aliases: []string{"tf"},
	Short:   "Manage terraform operations.",
	Long:    `Manage terraform operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("terraform called")
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}
