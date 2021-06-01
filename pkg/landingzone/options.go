//
// Rover - Landing zone & launchpad options
// * Used to hold the large number of params and vars for landing zone operations in a single place
// * Ben C, May 2021
//

package landingzone

import (
	"fmt"
	"os"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/spf13/cobra"
)

// Options holds all the settings for a langingzone or launchpad operation
// It's populated first by calling NewOptionsFromCmd, then the RunCmd func sets Subscription & Identity fields
type Options struct {
	LaunchPadMode      bool
	ConfigPath         string
	SourcePath         string
	Level              int
	CafEnvironment     string
	StateName          string
	Workspace          string
	TargetSubscription string
	StateSubscription  string
	Impersonate        bool
	OutPath            string
	Subscription       azure.Subscription
	Identity           azure.Identity
}

// NewOptionsFromCmd builds a Config from command flags
func NewOptionsFromCmd(cmd *cobra.Command) Options {
	launchPadMode := false
	if cmd.Parent().Name() == "launchpad" {
		launchPadMode = true
	}

	configPath, _ := cmd.Flags().GetString("config-path")
	sourcePath, _ := cmd.Flags().GetString("source")
	level, _ := cmd.Flags().GetInt("level")
	env, _ := cmd.Flags().GetString("environment")
	stateName, _ := cmd.Flags().GetString("statename")
	ws, _ := cmd.Flags().GetString("workspace")
	stateSub, _ := cmd.Flags().GetString("state-sub")
	targetSub, _ := cmd.Flags().GetString("target-sub")

	// This is a 'just in case' default, it will be changed later
	outPath, err := os.UserHomeDir()
	cobra.CheckErr(err)

	o := Options{
		LaunchPadMode:      launchPadMode,
		ConfigPath:         configPath,
		SourcePath:         sourcePath,
		Level:              level,
		CafEnvironment:     env,
		StateName:          stateName,
		Workspace:          ws,
		TargetSubscription: targetSub,
		StateSubscription:  stateSub,
		OutPath:            outPath,
	}

	return o
}

// LevelString returns the level as formated string
// This should be used rather than accessing level directly
func (o Options) LevelString() string {
	return fmt.Sprintf("level%d", o.Level)
}
