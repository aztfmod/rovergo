package test

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

func Test_Execute_Group_Deploy_Command(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)
	copyCommandYamlToRoverHome(roverHome, "group_no_commands.yml", "commands.yml")
	console.DebugEnabled = true
	testDataPath := "../../test/testdata"
	fmt.Println(testDataPath)

	//act
	custom.InitializeCustomCommandsAndGroups()
	deployAction := actions.ActionMap["deploy"]

	deployOptions := &cobra.Command{}
	deployOptions.Flags().String("config-dir", testDataPath+"/configs/level0/launchpad", "")
	deployOptions.Flags().String("source", os.Getenv("HOME")+"/.rover/caf-terraform-landingzones", "")
	deployOptions.Flags().String("level", "level0", "")
	deployOptions.Flags().String("environment", "test", "")
	deployOptions.Flags().Bool("launchpad", true, "")
	optionsList := landingzone.BuildOptions(deployOptions)

	err := deployAction.Execute(&optionsList[0])

	//assert
	assert.Equal(t, nil, err)

	t.Cleanup(func() {
		removeCommandYamlFromHomeDir(roverHome)
	})
}

func getTestHarnessPath(rootPath string) string {
	testPath := filepath.Join(rootPath, "test")
	testDataPath := filepath.Join(testPath, "testdata")
	return filepath.Join(testDataPath, "custom_commands")
}

func getProjectRootDir(currentWorkingDirectory string) string {
	pgk := filepath.Dir(currentWorkingDirectory)
	root := filepath.Dir(pgk)
	return root
}

func copyCommandYamlToRoverHome(roverHome, fileName string, target string) {
	currentWorkingDirectory, _ := os.Getwd()
	rootPath := getProjectRootDir(currentWorkingDirectory)
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
