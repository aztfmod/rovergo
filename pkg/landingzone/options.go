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
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
)

// Options holds all the settings for a langingzone or launchpad operation
// It's populated first by calling NewOptionsFromCmd, then the RunCmd func sets Subscription & Identity fields
// TODO: This probably needs a better name like `Operation` or something
type Options struct {
	LaunchPadMode      bool
	ConfigPath         string
	SourcePath         string
	Level              string
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

const cafLaunchPadDir = "/caf_launchpad"
const cafLandingzoneDir = "/caf_solution"

// NewOptionsFromCmd builds a Config from command flags
func NewOptionsFromCmd(cmd *cobra.Command) Options {
	configPath, _ := cmd.Flags().GetString("config-path")
	sourcePath, _ := cmd.Flags().GetString("source")
	level, _ := cmd.Flags().GetString("level")
	env, _ := cmd.Flags().GetString("environment")
	stateName, _ := cmd.Flags().GetString("statename")
	ws, _ := cmd.Flags().GetString("workspace")
	stateSub, _ := cmd.Flags().GetString("state-sub")
	targetSub, _ := cmd.Flags().GetString("target-sub")

	// Handle the launchpad mode special case
	launchPadMode := false
	if cmd.Parent().Name() == "launchpad" {
		launchPadMode = true
		// TODO: Maybe we have to remove this assumption and add --level flag to the `launchpad` cmd 😥
		level = "level0"
	}

	// This is a 'just in case' default, it will be changed later, when initializeCAF is called
	outPath, err := utils.GetHomeDirectory()
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
	}

	// Safely set the paths up
	o.SetSourcePath(sourcePath)
	o.SetConfigPath(configPath)

	return o
}

// LevelString returns the level as formated string
// This should be used rather than accessing level directly
// func (o Options) LevelString() string {
// 	return fmt.Sprintf("level%d", o.Level)
// }

// SetSourcePath ensures the source path is correct and absolute
func (o *Options) SetSourcePath(sourcePath string) {
	if strings.HasSuffix(sourcePath, cafLaunchPadDir) || strings.HasSuffix(sourcePath, cafLandingzoneDir) {
		cobra.CheckErr(fmt.Sprintf("source should not include %s or %s", cafLandingzoneDir, cafLaunchPadDir))
	}

	// Convert to absolute paths as a precaution
	sourcePath, err := filepath.Abs(sourcePath)
	cobra.CheckErr(err)

	if o.LaunchPadMode {
		o.SourcePath = path.Join(sourcePath, cafLaunchPadDir)
	} else {
		o.SourcePath = path.Join(sourcePath, cafLandingzoneDir)
	}

	// Try to ensure sourcepath is "good", i.e. exists & has some some terraform in it
	_, err = os.Stat(o.SourcePath)
	if err != nil {
		console.Errorf("Unable to open source directory: %s\n", o.SourcePath)
		cobra.CheckErr("Source directory must exist for rover to run")
	}
	foundTfFiles := false
	sourceFiles, err := os.ReadDir(o.SourcePath)
	cobra.CheckErr(err)
	for _, file := range sourceFiles {
		if strings.HasSuffix(file.Name(), ".tf") {
			foundTfFiles = true
			break
		}
	}
	if !foundTfFiles {
		console.Errorf("No terraform was found in source directory: %s\n", o.SourcePath)
		cobra.CheckErr("Rover execution has ended")
	}
}

func (o *Options) SetConfigPath(configPath string) {
	configPath, err := filepath.Abs(configPath)
	cobra.CheckErr(err)
	o.ConfigPath = configPath

	_, err = os.Stat(o.ConfigPath)
	if err != nil {
		console.Errorf("Unable to open config directory: %s\n", o.ConfigPath)
		cobra.CheckErr("Config directory must exist for rover to run")
	}
}
