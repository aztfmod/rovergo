//go:build (integration && ignore) || !unit
// +build integration,ignore !unit

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

func Test_Execute_RoverTest_OnExamples(t *testing.T) {
	//arrange
	roverHome := "/tmp"
	removeCommandYamlFromCWD()
	rover.SetHomeDirectory(roverHome)

	console.DebugEnabled = true
	testDataPath := "../../test/testdata"
	exampleTestPath := "../../examples/tests"
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

	copyCommandYamlToRoverHome(roverHome, "valid_group.yml", "commands.yml")
	console.DebugEnabled = true

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
