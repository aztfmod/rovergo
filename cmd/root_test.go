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
	assert.Equal(t, found.Short, "Perform a terraform init and no other action")
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
	assert.Equal(t, found.Short, "Perform a terraform plan")
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
	assert.Equal(t, found.Short, "Perform a terraform plan & apply")
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
	assert.Equal(t, found.Short, "Perform a terraform destroy")
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
	assert.Equal(t, found.Short, "Perform a terraform validate")
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
	assert.Equal(t, found.Short, "Perform a terraform format")
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
	assert.Equal(t, found.Short, "formats configuration files using terraform")
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
	assert.Equal(t, found.Short, "checks configuration files using terraform")
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
	assert.Equal(t, found.Short, "workflow used for CI/CD\n                  - plan\n                  - fmt\n                  - apply\n                  - validate\n                  - destroy\n")
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

	BuildSubCommandsFromActionMap()
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
