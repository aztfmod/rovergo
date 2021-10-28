//go:build unit
// +build unit

package custom

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/builtin/actions"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Valid_Commands_File_Exists(t *testing.T) {
	copyCommandYamlToCWD("_default.yml", "commands.yml")

	console.DebugEnabled = true

	actions, _ := LoadCustomCommandsAndGroups()

	assert.NotEmpty(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromCWD()
	})
}

func Test_CommandsFile_Not_In_CWD_And_Not_In_Rover_Home(t *testing.T) {
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	removeCommandYamlFromHomeDir(roverHome)
	rover.SetHomeDirectory(roverHome)
	console.DebugEnabled = true

	actions, err := LoadCustomCommandsAndGroups()

	assert.Error(t, err)
	assert.Empty(t, actions)
}

func Test_CommandsFile_FullExtension_Not_In_CWD_And_In_Rover_Home(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "_default.yml", "commands.yaml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.Nil(t, err)
	assert.NotEmpty(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_CommandsFile_Not_In_CWD_And_In_Rover_Home(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "_default.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.Nil(t, err)
	assert.NotEmpty(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Empty_CommandsFile_In_Rover_Home(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "empty.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.NotNil(t, err)
	assert.Empty(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_InvalidYaml_In_CommandsFile_In_Rover_Home(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "invalid_yaml.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.EqualError(t, err, "invalid yaml in /tmp/commands.yml. Internal Error:yaml: unmarshal errors:\n  line 3: field not valid --- \" not found in type custom.Command")
	assert.Nil(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Custom_Command_Name_Collision_With_Built_In_Command(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "builtin_command_collision.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.EqualError(t, err, "invalid custom command (plan). Custom command (plan) cannot be used as it is a builtin command")
	assert.Nil(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Group_Name_Collision_With_Built_In_Command(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "group_name_collision.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.Error(t, err, "invalid group name (plan). (plan) cannot be used as it is a builtin command")
	assert.Nil(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Group_With_Invalid_Command(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "group_invalid_command.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.Error(t, err, "invalid group name (foo). (foo) must be a valid built in command or a custom command")
	assert.Nil(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Groups_With_Custom_Commands_Are_Allowed(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "valid_group.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.Nil(t, err)
	assert.NotEmpty(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Groups_With_EmptyCommands_Are_NotAllowed(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "group_empty_commands.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.EqualError(t, err, "invalid group (deploy). A group must have at least one command")
	assert.Nil(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_CommandsYaml_WithGroupsSection_NoCustomCommandSection_Is_Allowed(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "group_no_commands.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	actions, err := LoadCustomCommandsAndGroups()

	//assert
	assert.Nil(t, err)
	assert.NotNil(t, actions)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_InitilizeCustomCommands_ActionMap_Contains_CustomCommand(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "valid_group.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	InitializeCustomCommandsAndGroups()
	//assert

	exists := contains(actions.ActionMap, "format")
	assert.True(t, exists)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_InitilizeCustomCommands_ActionMap_Contains_Group(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "valid_group.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	InitializeCustomCommandsAndGroups()
	//assert

	exists := contains(actions.ActionMap, "deploy")
	assert.True(t, exists)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_InitilizeCustomCommands_Group_Contains_Expected_Commands(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "valid_group.yml", "commands.yml")
	console.DebugEnabled = true

	//act
	InitializeCustomCommandsAndGroups()
	//assert

	deploy := actions.ActionMap["deploy"].(Action)
	assert.Equal(t, 3, len(deploy.Commands)) // 3 commands are in the test harness file valid_groups.yml

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func Test_Execute_Test(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	console.DebugEnabled = true
	testDataPath := "../../test/testdata"
	exampleTestPath := "../../examples/tests"
	fmt.Println(testDataPath)

	testOptions := &cobra.Command{}
	testOptions.Flags().String("config-dir", testDataPath+"/configs/level0/launchpad", "")
	testOptions.Flags().String("test-source", exampleTestPath, "")
	testOptions.Flags().String("level", "level0", "")
	testOptions.Flags().String("environment", "test", "")
	testOptions.Flags().Bool("launchpad", true, "")
	sub, _ := azure.GetSubscription()
	testOptions.Flags().String("state-sub", sub.ID, "")
	testOptions.Flags().String("statename", "caf_launchpad", "")

	optionsList := landingzone.BuildOptions(testOptions)

	//act
	InitializeCustomCommandsAndGroups()
	testAction := actions.ActionMap["test"]

	//assert
	assert.Equal(t, nil, testAction.Execute(&optionsList[0]))

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func getTestHarnessPath(rootPath string) string {
	testPath := filepath.Join(rootPath, "test")
	testDataPath := filepath.Join(testPath, "testdata")
	return filepath.Join(testDataPath, "custom_commands")
}

func copyCommandYamlToCWD(fileName string, target string) {
	currentWorkingDirectory, _ := os.Getwd()
	rootPath := GetProjectRootDir(currentWorkingDirectory)
	testHarnessPath := getTestHarnessPath(rootPath)
	sourcePath := filepath.Join(testHarnessPath, fileName)
	destinationPath := filepath.Join(currentWorkingDirectory, target)
	utils.CopyFile(sourcePath, destinationPath)
}

func copyCommandYamlToRoverHome(roverHome, fileName string, target string) {
	currentWorkingDirectory, _ := os.Getwd()
	rootPath := GetProjectRootDir(currentWorkingDirectory)
	testHarnessPath := getTestHarnessPath(rootPath)
	sourcePath := filepath.Join(testHarnessPath, fileName)
	destinationPath := filepath.Join(roverHome, target)
	utils.CopyFile(sourcePath, destinationPath)
}

func removeCommandYamlFromCWD() {
	fileName := "commands.yml"
	currentWorkingDirectory, _ := os.Getwd()
	filePath := filepath.Join(currentWorkingDirectory, fileName)
	e := os.Remove(filePath)
	if e != nil {
		_ = fmt.Errorf("Error removing test harness command.yml - %s", e)
	}
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
	pgk := filepath.Dir(currentWorkingDirectory)
	root := filepath.Dir(pgk)
	return root
}
