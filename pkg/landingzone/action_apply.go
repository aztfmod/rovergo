package landingzone

import (
	"context"
	"fmt"
	"os"

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
			tfexec:             nil,
			ActionBase: ActionBase{
				name:        "apply",
				description: "Perform a terraform plan & apply",
			},
		},
	}
}

func (a *ApplyAction) Execute(o *Options) error {

	planAction := NewPlanAction()
	_ = planAction.Execute(o)

	if !planAction.hasChanges {
		console.Success("Plan resulted in no changes, apply will be skipped")
		return nil
	}

	console.Info("Carrying out Terraform apply")

	// We pull these from the previous action, which is a special case
	// But saves us re-running prepareTerraformCAF, or
	a.tfexec = planAction.tfexec
	a.launchPadStorageID = planAction.launchPadStorageID

	planFile := fmt.Sprintf("%s/%s.tfplan", o.OutPath, o.StateName)
	stateFile := fmt.Sprintf("%s/%s.tfstate", o.OutPath, o.StateName)

	// Build apply options, with plan file and state out
	applyOptions := []tfexec.ApplyOption{
		tfexec.DirOrPlan(planFile),
		tfexec.StateOut(stateFile),
		tfexec.Parallelism(terraformParallelism),
	}

	console.StartSpinner()
	err := a.tfexec.Apply(context.Background(), applyOptions...)
	console.StopSpinner()
	cobra.CheckErr(err)

	// Special case for post launchpad deployment
	newStorageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	cobra.CheckErr(err)
	if o.LaunchPadMode && a.launchPadStorageID != newStorageID {
		console.Info("Detected the launchpad infrastructure has been deployed or updated")

		stateFileName := o.OutPath + "/" + o.StateName + ".tfstate"
		azure.UploadFileToBlob(newStorageID, o.Workspace, o.StateName+".tfstate", stateFileName)
		console.Info("Uploading state from launchpad process to Azure storage")
		os.Remove(stateFileName)
	}

	console.Success("Apply was successful")

	return nil
}
