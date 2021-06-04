//
// Rover - support for landingzone / launchpad cobra cmds
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
)

// RunFunc is the shared cobra run function for `landingzone run` and `launchpad run`
func RunFunc(cmd *cobra.Command, args []string) {
	actionStr, _ := cmd.Flags().GetString("action")
	action, err := NewAction(actionStr)
	cobra.CheckErr(err)
	console.Infof("Starting '%s' command with action '%s'\n", cmd.CommandPath(), actionStr)

	// Build config from command flags
	opt := NewOptionsFromCmd(cmd)

	opt.Execute(action)
}

// SetSharedFlags configures command flags for both landingzone and launchpad commands
func SetSharedFlags(cmd *cobra.Command) {
	cmd.PersistentFlags()
	cmd.Flags().StringP("source", "s", "", "Path to source of landingzone (required)")
	cmd.Flags().StringP("config-path", "c", "", "Configuration vars directory (required)")
	cmd.Flags().StringP("environment", "e", "sandpit", "Name of CAF environment")
	cmd.Flags().StringP("workspace", "w", "tfstate", "Name of workspace")
	cmd.Flags().StringP("statename", "n", "", "Name for state and plan files, (default landingzone source dir name)")
	cmd.Flags().String("state-sub", "", "Azure subscription ID where state is held")
	cmd.Flags().String("target-sub", "", "Azure subscription ID to operate on")
	AddActionFlag(cmd)
	cmd.Flags().SortFlags = true

	// Level command not on launchpad cmd as we fix it to "level0"
	if cmd.Parent().Name() != "launchpad" {
		cmd.Flags().StringP("level", "l", "level1", "Level")
	}

	_ = cobra.MarkFlagRequired(cmd.Flags(), "source")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "config-path")
}
