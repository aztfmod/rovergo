//
// Rover - Top level cd (continuous deployment) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Manage CD operations.",
	Long:  `Manage CD operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cd called")
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
