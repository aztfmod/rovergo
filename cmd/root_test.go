//go:build unit
// +build unit

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aztfmod/rover/pkg/builtin/actions"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/custom"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/aztfmod/rover/pkg/symphony"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const testDataPath = "../test/testdata"

func Test_Rover_Standalone_Apply_Launchpad(t *testing.T) {
	console.DebugEnabled = true

	testCmd := &cobra.Command{
		Use: "apply",
	}
	testCmd.Flags().String("config-dir", testDataPath+"/configs/level0/launchpad", "")
	testCmd.Flags().String("source", testDataPath+"/caf-terraform-landingzones", "")
	testCmd.Flags().String("level", "level0", "")
	testCmd.Flags().Bool("launchpad", true, "")

	optionsList := landingzone.BuildOptions(testCmd)

	configPath, err := filepath.Abs(testDataPath + "/configs/level0/launchpad")
	if err != nil {
		t.Fail()
	}

	sourcePath, err := filepath.Abs(testDataPath + "/caf-terraform-landingzones/caf_launchpad")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, optionsList[0].ConfigPath, configPath)
	assert.Equal(t, optionsList[0].SourcePath, sourcePath)
	assert.Equal(t, optionsList[0].CafEnvironment, "sandpit")
	assert.Equal(t, optionsList[0].StateName, "caf_launchpad")
	assert.Equal(t, optionsList[0].Workspace, "tfstate")
	assert.Equal(t, optionsList[0].DryRun, false)
	assert.Equal(t, optionsList[0].TargetSubscription, "")
	assert.Equal(t, optionsList[0].StateSubscription, "")
	assert.Equal(t, optionsList[0].LaunchPadMode, true)

	getActionMap()
	action := actions.ActionMap["mock"]
	_ = action.Execute(&optionsList[0])
}

func Test_Builtin_Init_Command(t *testing.T) {
	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "init" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "init")
	assert.Equal(t, found.Short, "[Builtin command]\tPerform a terraform init and no other action")
	assert.Equal(t, found.Long, "")
}

func Test_Builtin_Plan_Command(t *testing.T) {
	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "plan" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "plan")
	assert.Equal(t, found.Short, "[Builtin command]\tPerform a terraform plan")
	assert.Equal(t, found.Long, "")
}

func Test_Builtin_Apply_Command(t *testing.T) {
	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "apply" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "apply")
	assert.Equal(t, found.Short, "[Builtin command]\tPerform a terraform plan & apply")
	assert.Equal(t, found.Long, "")
}

func Test_Builtin_Destroy_Command(t *testing.T) {
	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "destroy" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "destroy")
	assert.Equal(t, found.Short, "[Builtin command]\tPerform a terraform destroy")
	assert.Equal(t, found.Long, "")
}

func Test_Builtin_Validate_Command(t *testing.T) {
	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "validate" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "validate")
	assert.Equal(t, found.Short, "[Builtin command]\tPerform a terraform validate")
	assert.Equal(t, found.Long, "")
}

func Test_Builtin_Fmt_Command(t *testing.T) {
	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "fmt" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "fmt")
	assert.Equal(t, found.Short, "[Builtin command]\tPerform a terraform format")
	assert.Equal(t, found.Long, "")
}

func Test_Custom_Format_Command(t *testing.T) {
	roverHome := "/tmp"
	rover.SetHomeDirectory(roverHome)

	copyCommandYamlToRoverHome(roverHome, "_default.yml", "commands.yml")

	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "format" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "format")
	assert.Equal(t, found.Short, "[Custom command]\tPerform terraform with list,write parameters")
	assert.Equal(t, found.Long, "")

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Custom_Check_Command(t *testing.T) {
	roverHome := "/tmp"
	rover.SetHomeDirectory(roverHome)

	copyCommandYamlToRoverHome(roverHome, "_default.yml", "commands.yml")

	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "check" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "check")
	assert.Equal(t, found.Short, "[Custom command]\tPerform terraform with list,write parameters")
	assert.Equal(t, found.Long, "")

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Group_Deploy_Command(t *testing.T) {
	roverHome := "/tmp"
	rover.SetHomeDirectory(roverHome)

	copyCommandYamlToRoverHome(roverHome, "_default.yml", "commands.yml")

	console.DebugEnabled = true

	getActionMap()

	var found *cobra.Command

	allCommands := rootCmd.Commands()
	for _, cmd := range allCommands {
		if cmd.Use == "deploy" {
			found = cmd
		}
	}

	if found == nil {
		t.Fail()
	}

	assert.Equal(t, found.Use, "deploy")
	assert.Equal(t, found.Short, "[Group command]\tPerform plan,fmt,apply,validate,destroy commands sequentially")
	assert.Equal(t, found.Long, "")

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func getActionMap() {
	err := custom.InitializeCustomCommandsAndGroups()
	if err != nil {
		console.Warningf("No custom command or group found in the current directory or rover home directory\n")
	}

	actions.ActionMap["mock"] = newMockAction()

	// Dynamically build sub-commands from list of actions
	for key, action := range actions.ActionMap {
		actionSubCmd := &cobra.Command{
			Use:   key,
			Short: fmt.Sprintf("[%s command]\t%s", action.GetType(), action.GetDescription()),
			PreRun: func(cmd *cobra.Command, args []string) {
			},
			Run: func(cmd *cobra.Command, args []string) {
				// NOTE: We CAN NOT use the action variable from the loop above as it's not bound at runtime
				// Dynamically building our commands has some limitations, instead we need to use the cmd name & the map
				action = actions.ActionMap[cmd.Name()]

				configFile, _ := cmd.Flags().GetString("config-file")
				configPath, _ := cmd.Flags().GetString("config-dir")

				// Handle the user trying to use both configPath and configFile or neither!
				if configPath == "" && configFile == "" {
					_ = cmd.Help()
					os.Exit(0)
				}
				if configPath != "" && configFile != "" {
					cobra.CheckErr("--config-file and --config-dir options must not be combined, specify only one")
				}

				var optionsList []landingzone.Options
				// Handle symphony mode where config file and level is passed, this will return optionsList with MANY items
				if configFile != "" {
					// Depending on if we're running single or mult-level this will return one or many options
					optionsList = symphony.BuildOptions(cmd)
				}

				// Handle CLI or standalone mode, this will return optionsList with a single item
				if configPath != "" {
					optionsList = landingzone.BuildOptions(cmd)
				}

				for _, options := range optionsList {
					// Now start the action execution...
					// If an error occurs, depend on downstream code to log messages
					console.Infof("Executing action %s for %s\n", action.GetName(), options.StateName)
					err = action.Execute(&options)
					if err != nil {
						cobra.CheckErr(err)
					}
				}

				console.Success("Rover has finished")
				os.Exit(0)
			},
		}

		actionSubCmd.Flags().StringP("source", "s", "", "Path to source of landingzone")
		actionSubCmd.Flags().StringP("config-file", "c", "", "Configuration file, you must supply this or config-dir")
		actionSubCmd.Flags().StringP("config-dir", "v", "", "Configuration directory, you must supply this or config-file")
		actionSubCmd.Flags().StringP("environment", "e", "", "Name of CAF environment")
		actionSubCmd.Flags().StringP("workspace", "w", "", "Name of workspace")
		actionSubCmd.Flags().StringP("statename", "n", "", "Name for state and plan files")
		actionSubCmd.Flags().String("state-sub", "", "Azure subscription ID where state is held")
		actionSubCmd.Flags().String("target-sub", "", "Azure subscription ID to operate on")
		actionSubCmd.Flags().Bool("launchpad", false, "Run in launchpad mode, i.e. level 0")
		actionSubCmd.Flags().StringP("level", "l", "", "CAF landingzone level name, default is all levels")
		actionSubCmd.Flags().BoolP("dry-run", "d", false, "Execute a dry run where no actions will be executed")
		actionSubCmd.Flags().SortFlags = true

		// Stuff it under the parent root command
		rootCmd.AddCommand(actionSubCmd)
	}
}

func newMockAction() *landingzone.MockAction {
	return &landingzone.MockAction{
		ActionBase: landingzone.ActionBase{
			Name:        "mock",
			Description: "do nothing",
		},
	}
}

func getTestHarnessPath(rootPath string) string {
	testPath := filepath.Join(rootPath, "test")
	testDataPath := filepath.Join(testPath, "testdata")
	return filepath.Join(testDataPath, "custom_commands")
}

func copyCommandYamlToRoverHome(roverHome, fileName string, target string) {
	currentWorkingDirectory, _ := os.Getwd()
	rootPath := GetProjectRootDir(currentWorkingDirectory)
	testHarnessPath := getTestHarnessPath(rootPath)
	sourcePath := filepath.Join(testHarnessPath, fileName)
	destinationPath := filepath.Join(roverHome, target)
	utils.CopyFile(sourcePath, destinationPath)
}

func removeCommandYamlFromHomeDir(homeDir string) {
	fileNames := [2]string{"commands.yml", "commands.yaml"}
	for _, fileName := range fileNames {
		filePath := filepath.Join(homeDir, fileName)
		e := os.Remove(filePath)
		if e != nil {
			_ = fmt.Errorf("Error removing test harness %s - %s", fileName, e)
		}
	}
}

func GetProjectRootDir(currentWorkingDirectory string) string {
	return filepath.Dir(currentWorkingDirectory)
}
