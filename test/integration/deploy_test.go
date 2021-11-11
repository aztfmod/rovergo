//go:build integration && !unit
// +build integration,!unit

package test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	rand.Seed(time.Now().UnixNano())
	environmentName := fmt.Sprintf("ci_%d", rand.Intn(100))

	testDataPath := "../../test/testdata"
	fmt.Println(testDataPath)

	//act
	err := custom.InitializeCustomCommandsAndGroups()
	if err != nil {
		t.Errorf("Error initializing custom commands and groups - %s", err)
	}

	deployAction := actions.ActionMap["deploy"]

	deployOptions := &cobra.Command{}
	deployOptions.Flags().String("config-dir", testDataPath+"/configs/level0/launchpad", "")
	deployOptions.Flags().String("source", os.Getenv("HOME")+"/.rover/caf-terraform-landingzones", "")
	deployOptions.Flags().String("level", "level0", "")
	deployOptions.Flags().String("environment", environmentName, "")
	deployOptions.Flags().Bool("launchpad", true, "")
	optionsList := landingzone.BuildOptions(deployOptions)

	err = deployAction.Execute(&optionsList[0])

	//assert
	assert.Equal(t, nil, err)

	t.Cleanup(func() {
		actions.ActionMap["destroy"].Execute(&optionsList[0])

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
	err := utils.CopyFile(sourcePath, destinationPath)
	if err != nil {
		_ = fmt.Errorf("Error copying test harness %s - %s", fileName, err)
	}
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
