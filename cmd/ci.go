//
// Rover - Top level ci (continuous integration) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Manage CI operations.",
	Long:  `Manage CI operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ci called")
	},
}

func init() {
	rootCmd.AddCommand(ciCmd)
}
