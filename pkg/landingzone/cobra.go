//
// Rover - support for landingzone / launchpad cobra cmds
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"path/filepath"

	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
)

const defaultWorkspace = "tfstate"
const defaultEnv = "sandpit"

// BuildOptions parses the CLI params and flags and return a single Option
// Note. it's returned as an single item array for symmetry with symphony.BuildOptions
func BuildOptions(cmd *cobra.Command) []Options {
	launchPadMode, _ := cmd.Flags().GetBool("launchpad")
	configPath, _ := cmd.Flags().GetString("config-dir")
	sourcePath, _ := cmd.Flags().GetString("source")
	level, _ := cmd.Flags().GetString("level")
	env, _ := cmd.Flags().GetString("environment")
	stateName, _ := cmd.Flags().GetString("statename")
	ws, _ := cmd.Flags().GetString("workspace")
	stateSub, _ := cmd.Flags().GetString("state-sub")
	targetSub, _ := cmd.Flags().GetString("target-sub")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// This is a 'just in case' default, it will be changed later, when initializeCAF is called
	outPath, err := utils.GetRoverDirectory()
	cobra.CheckErr(err)

	// Normally cobra would provide automatic defaults but we are using it in a weird way
	if level == "" {
		cobra.CheckErr("--level must be supplied when not using a config file")
	}
	if sourcePath == "" {
		cobra.CheckErr("--source option must be supplied when not using a config file")
	}
	if ws == "" {
		ws = defaultWorkspace
	}
	if env == "" {
		env = defaultEnv
	}

	opt := Options{
		LaunchPadMode:      launchPadMode,
		Level:              level,
		CafEnvironment:     env,
		StateName:          stateName,
		Workspace:          ws,
		TargetSubscription: targetSub,
		StateSubscription:  stateSub,
		OutPath:            outPath,
		DryRun:             dryRun,
	}

	// Safely set the paths up
	opt.SetSourcePath(sourcePath)
	opt.SetConfigPath(configPath)

	// Default state & plan name is taken from the base name of the landingzone source dir
	if stateName == "" {
		stateName = filepath.Base(opt.SourcePath)
		if stateName == "/" || stateName == "." {
			cobra.CheckErr("Error --source should be a directory path")
		}
		// Update the StateName, we have to do this after SetSourcePath is called
		opt.StateName = stateName
	}

	return []Options{opt}
}
