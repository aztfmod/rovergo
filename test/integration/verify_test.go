package test

//TODO: Uncomment and fix once integration test pattern is finalized.

// //go:build (integration && ignore) || !unit
// // +build integration,ignore !unit

// import (
// 	"testing"

// 	"github.com/aztfmod/rover/pkg/azure"
// 	"github.com/aztfmod/rover/pkg/builtin/actions"
// 	"github.com/aztfmod/rover/pkg/console"
// 	"github.com/aztfmod/rover/pkg/custom"
// 	"github.com/aztfmod/rover/pkg/landingzone"
// 	"github.com/aztfmod/rover/pkg/rover"
// 	"github.com/spf13/cobra"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_Execute_RoverTest_OnExamples(t *testing.T) {
// 	//arrange
// 	roverHome := "/tmp"
// 	custom.removeCommandYamlFromCWD()
// 	rover.SetHomeDirectory(roverHome)

// 	console.DebugEnabled = true
// 	testDataPath := "../../test/testdata"
// 	exampleTestPath := "../../examples/tests"
// 	testOptions := &cobra.Command{}
// 	testOptions.Flags().String("config-dir", testDataPath+"/configs/level0/launchpad", "")
// 	testOptions.Flags().String("test-source", exampleTestPath, "")
// 	testOptions.Flags().String("level", "level0", "")
// 	testOptions.Flags().String("environment", "test", "")
// 	testOptions.Flags().Bool("launchpad", true, "")
// 	sub, _ := azure.GetSubscription()
// 	testOptions.Flags().String("state-sub", sub.ID, "")
// 	testOptions.Flags().String("statename", "caf_launchpad", "")

// 	optionsList := landingzone.BuildOptions(testOptions)

// 	//act
// 	custom.InitializeCustomCommandsAndGroups()
// 	testAction := actions.ActionMap["test"]

// 	//assert
// 	assert.Equal(t, nil, testAction.Execute(&optionsList[0]))

// 	custom.copyCommandYamlToRoverHome(roverHome, "valid_group.yml", "commands.yml")
// 	console.DebugEnabled = true

// }
