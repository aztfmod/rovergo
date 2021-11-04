package landingzone

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type ApplyAction struct {
	TerraformAction
}

func NewApplyAction() *ApplyAction {
	return &ApplyAction{
		TerraformAction: TerraformAction{
			launchPadStorageID: "",
			ActionBase: ActionBase{
				Name:        "apply",
				Type:        BuiltinCommand,
				Description: "Perform a terraform plan & apply",
			},
		},
	}
}

func (a *ApplyAction) Execute(o *Options) error {
	tf, err := a.prepareTerraformCAF(o)
	if err != nil {
		return err
	}

	planFile := path.Join(o.DataDir, fmt.Sprintf("%s.tfplan", o.StateName))
	stateFile := path.Join(o.DataDir, fmt.Sprintf("%s.tfstate", o.StateName))
	console.Infof("Apply will use plan file %s\n", planFile)

	// Build apply options, with plan file and state out
	applyOptions := []tfexec.ApplyOption{
		tfexec.DirOrPlan(planFile),
		tfexec.StateOut(stateFile),
		tfexec.Parallelism(terraformParallelism),
	}

	console.StartSpinner()
	err = tf.Apply(context.Background(), applyOptions...)
	console.StopSpinner()
	cobra.CheckErr(err)

	// Special case for post launchpad deployment
	newStorageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	cobra.CheckErr(err)
	if o.LaunchPadMode && a.launchPadStorageID != newStorageID {
		console.Info("Detected the launchpad infrastructure has been deployed or updated")

		stateFileName := o.DataDir + "/" + o.StateName + ".tfstate"
		err := azure.UploadFileToBlob(newStorageID, o.Workspace, o.StateName+".tfstate", stateFileName)
		cobra.CheckErr(err)
		console.Info("Uploading state from launchpad process to Azure storage")
		os.Remove(stateFileName)

		// Why re-init with remote this straight after?
		// Otherwise we aren't tracking state at all, state will be uploaded to Azure but we won't use it
		err = o.runRemoteInit(tf, newStorageID)
		cobra.CheckErr(err)
	}

	console.Success("Apply was successful")
	console.Infof("Removing plan file: %s\n", planFile)
	_ = os.Remove(planFile)

	return nil
}
