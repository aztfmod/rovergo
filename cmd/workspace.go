//
// Rover - Top level workspace command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage workspace operations.",
	Long:    `Manage workspace operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspace called")
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
}
