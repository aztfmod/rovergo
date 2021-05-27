//
// Rover - Config command
// * Does nothing, branches down to sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Access to configuration related sub-commands, such as 'auth'.",
	Long:    `Access to configuration related sub-commands, such as 'auth'.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
