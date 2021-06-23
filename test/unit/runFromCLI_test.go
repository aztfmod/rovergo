package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aztfmod/rover/cmd"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/custom"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Rover_Standalone_Apply_Launchpad(t *testing.T) {

	console.DebugEnabled = true

	testCmd := &cobra.Command{
		Use: "apply",
	}
	testCmd.Flags().String("config-dir", "../testdata/configs/level0/launchpad", "")
	testCmd.Flags().String("source", "../testdata/caf-terraform-landingzones", "")
	testCmd.Flags().String("level", "level0", "")
	testCmd.Flags().Bool("launchpad", true, "")

	optionsList := landingzone.BuildOptions(testCmd)

	configPath, err := filepath.Abs("../testdata/configs/level0/launchpad")
	if err != nil {
		t.Fail()
	}

	sourcePath, err := filepath.Abs("../testdata/caf-terraform-landingzones/caf_launchpad")
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
	action := cmd.ActionMap["mock"]
	_ = action.Execute(&optionsList[0])

}

func getActionMap() {
	custActions, err := custom.FetchActions()
	if err != nil {
		console.Errorf("Failed %s", err)
		os.Exit(1)
	}
	for _, ca := range custActions {
		cmd.ActionMap[ca.GetName()] = ca
	}
	cmd.ActionMap["mock"] = NewMockAction()
}

func NewMockAction() *landingzone.MockAction {
	return &landingzone.MockAction{
		ActionBase: landingzone.ActionBase{
			Name:        "mock",
			Description: "do nothing",
		},
	}
}