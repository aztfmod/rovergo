package landingzone

import (
	"errors"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
)

type TestAction struct {
	ActionBase
}

func NewTestAction() *TestAction {
	return &TestAction{
		ActionBase{
			Name:        "Test",
			Type:        BuiltinCommand,
			Description: "Perform a Go Test",
		},
	}
}

func (ta *TestAction) Execute(o *Options) error {
	console.Info("Carrying out Rover Test")

	// download tfstate file

	// Lccate storage account id
	stateFileName := o.DataDir + "/" + o.StateName + ".tfstate"

	storageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	if err != nil {

		console.Errorf("No state storage account found for environment '%s' and level %s", o.CafEnvironment, o.Level)
		return errors.New("can't deploy a landing zone without a launchpad")
	} else {
		console.Infof("Located state storage account %s\n", storageID)
	}

	// doanload tfstate file
	err = azure.DownloadFileFromBlob(storageID, o.Workspace, o.StateName+".tfstate", stateFileName)
	cobra.CheckErr(err)

	// invoke go test

	// create junit test report.
	return nil
}
