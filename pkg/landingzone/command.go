//
// Rover - support for landingzone / launchpad cobra cmds
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"github.com/spf13/cobra"
)

// SetSharedFlags configures command flags for both landingzone and launchpad commands
func SetSharedFlags(cmd *cobra.Command) {
	cmd.PersistentFlags()
	cmd.PersistentFlags().StringP("source", "s", "", "Path to source of landingzone (required)")
	cmd.PersistentFlags().StringP("config-path", "c", "", "Configuration vars directory (required)")
	cmd.PersistentFlags().StringP("environment", "e", "sandpit", "Name of CAF environment")
	cmd.PersistentFlags().StringP("workspace", "w", "tfstate", "Name of workspace")
	cmd.PersistentFlags().StringP("statename", "n", "", "Name for state and plan files, (default landingzone source dir name)")
	cmd.PersistentFlags().String("state-sub", "", "Azure subscription ID where state is held")
	cmd.PersistentFlags().String("target-sub", "", "Azure subscription ID to operate on")
	cmd.Flags().SortFlags = true

	// Level flag not on launchpad cmd as we fix it to "level0"
	if cmd.Name() != "launchpad" {
		cmd.PersistentFlags().StringP("level", "l", "level1", "Level")
	}

	_ = cobra.MarkFlagRequired(cmd.Flags(), "source")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "config-path")
}
