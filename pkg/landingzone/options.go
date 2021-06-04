//
// Rover - Landing zone & launchpad options
// * Used to hold the large number of params and vars for landing zone operations in a single place
// * Ben C, May 2021
//

package landingzone

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	RunInit            bool
	Subscription       azure.Subscription
	Identity           azure.Identity
}

const cafLaunchPadDir = "/caf_launchpad"
const cafLandingzoneDir = "/caf_solution"

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

	// Default state & plan name is taken from the base name of the landingzone source dir
	if stateName == "" {
		stateName = filepath.Base(sourcePath)
		if stateName == "/" || stateName == "." {
			cobra.CheckErr("Error source-path should be a directory path")
		}
	}

	o := Options{
		LaunchPadMode:      launchPadMode,
		Level:              level,
		CafEnvironment:     env,
		StateName:          stateName,
		Workspace:          ws,
		TargetSubscription: targetSub,
		StateSubscription:  stateSub,
		OutPath:            outPath,
		RunInit:            true,
	}

	// Safely set the paths up
	o.SetSourcePath(sourcePath)
	o.SetConfigPath(configPath)

	return o
}

// LevelString returns the level as formated string
// This should be used rather than accessing level directly
func (o Options) LevelString() string {
	return fmt.Sprintf("level%d", o.Level)
}

// SetSourcePath ensures the source path is correct and absolute
func (o *Options) SetSourcePath(sourcePath string) {
	if strings.HasSuffix(sourcePath, cafLaunchPadDir) || strings.HasSuffix(sourcePath, cafLandingzoneDir) {
		cobra.CheckErr(fmt.Sprintf("source should not include %s or %s", cafLandingzoneDir, cafLaunchPadDir))
	}

	// TODO: Add validation that path exists and contains some .tf files

	// Convert to absolute paths as a precaution
	sourcePath, err := filepath.Abs(sourcePath)
	cobra.CheckErr(err)

	if o.LaunchPadMode {
		o.SourcePath = path.Join(sourcePath, cafLaunchPadDir)
	} else {
		o.SourcePath = path.Join(sourcePath, cafLandingzoneDir)
	}
}

func (o *Options) SetConfigPath(configPath string) {
	// TODO: Add validation that path exists and contains some .tfvar files

	configPath, err := filepath.Abs(configPath)
	cobra.CheckErr(err)
	o.ConfigPath = configPath
}
