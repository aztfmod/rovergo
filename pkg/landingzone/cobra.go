//
// Rover - support for landingzone / launchpad cobra cmds
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"path/filepath"

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
	stack, _ := cmd.Flags().GetString("stack")
	env, _ := cmd.Flags().GetString("environment")
	stateName, _ := cmd.Flags().GetString("statename")
	ws, _ := cmd.Flags().GetString("workspace")
	stateSub, _ := cmd.Flags().GetString("state-sub")
	targetSub, _ := cmd.Flags().GetString("target-sub")
	testpath, _ := cmd.Flags().GetString("test-source")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Normally cobra would provide automatic defaults but we are using it in a weird way
	if level == "" {
		cobra.CheckErr("--level must be supplied when not using a config file")
	}
	if sourcePath == "" && testpath == "" {
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
		Stack:              stack,
		CafEnvironment:     env,
		StateName:          stateName,
		Workspace:          ws,
		TargetSubscription: targetSub,
		StateSubscription:  stateSub,
		TestPath:           testpath,
		DryRun:             dryRun,
	}

	// Safely set the paths up

	opt.SetConfigPath(configPath)

	if testpath != "" {
		opt.SetTestPath(testpath)
	}

	if testpath == "" {
		opt.SetSourcePath(sourcePath)
		// Default state & plan name is taken from the base name of the landingzone source dir
		if stateName == "" {
			stateName = filepath.Base(opt.SourcePath)
			if stateName == "/" || stateName == "." {
				cobra.CheckErr("Error --source should be a directory path")
			}
			// Update the StateName, we have to do this after SetSourcePath is called
			opt.StateName = stateName
		}
	}

	err := opt.SetDataDir()
	cobra.CheckErr(err)

	return []Options{opt}
}
