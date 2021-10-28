//
// Rover - Landing zone & launchpad options
// * Used to hold the large number of params and vars for landing zone operations in a single place
// * Ben C, May 2021
//

package landingzone

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/spf13/cobra"
)

// Options holds all the settings for a langingzone or launchpad operation
// It's populated by NewOptionsFromCmd or from from YAML config, then the Execute func sets Subscription & Identity fields
type Options struct {
	LaunchPadMode      bool
	ConfigPath         string
	SourcePath         string
	TestPath           string
	Level              string
	Stack              string
	CafEnvironment     string
	StateName          string
	Workspace          string
	TargetSubscription string
	StateSubscription  string
	Impersonate        bool
	DataDir            string
	DryRun             bool
	Subscription       azure.Subscription
	Identity           azure.Identity
}

const cafLaunchPadDir = "/caf_launchpad"
const cafLandingzoneDir = "/caf_solution"

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

func (o *Options) SetTestPath(testPath string) {
	testPath, err := filepath.Abs(testPath)
	cobra.CheckErr(err)
	o.TestPath = testPath

	_, err = os.Stat(o.TestPath)
	if err != nil {
		console.Errorf("Unable to open Test directory: %s\n", o.TestPath)
		cobra.CheckErr("Test directory must exist for rover test to run")
	}
}

func (o *Options) SetDataDir() error {
	roverHome, err := rover.HomeDirectory()
	if err != nil {
		return err
	}
	// Build the data directory under the rover home
	// The hierarchy is: ~/.rover/workspace/level/statename
	o.DataDir = filepath.Join(roverHome, o.Workspace, o.Level, o.StateName)
	err = os.MkdirAll(o.DataDir, os.ModePerm)
	if err != nil {
		return err
	}
	console.Infof("Terraform data directory is set to: %s\n", o.DataDir)
	return nil
}

func (o *Options) Debug() {
	if !console.DebugEnabled {
		return
	}
	debugConf, _ := json.MarshalIndent(o, "", "  ")
	console.Debug(string(debugConf))
}
