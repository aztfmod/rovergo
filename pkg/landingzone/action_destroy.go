package landingzone

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type DestroyAction struct {
	TerraformAction
}

func NewDestroyAction() *DestroyAction {
	return &DestroyAction{
		TerraformAction: TerraformAction{
			launchPadStorageID: "",
			ActionBase: ActionBase{
				Name:        "destroy",
				Type:        BuiltinCommand,
				Description: "Perform a terraform destroy",
			},
		},
	}
}

func (a *DestroyAction) Execute(o *Options) error {
	tf, err := a.prepareTerraformCAF(o)
	if err != nil {
		return err
	}

	if o.DryRun {
		return nil
	}

	stateFileName := path.Join(o.DataDir, fmt.Sprintf("%s.tfstate", o.StateName))

	// Build apply options, with plan file and state out
	destroyOptions := []tfexec.DestroyOption{
		tfexec.Parallelism(terraformParallelism),
		tfexec.Refresh(false),
	}

	// We need to do all sorts of extra shenanigans for launchPadMode
	if o.LaunchPadMode {
		console.Warning("WARNING! You are destroying the launchpad!")
		if a.launchPadStorageID == "" {
			console.Error("Looks like this launchpad has already been deleted, bye!")
			cobra.CheckErr("Destroy was aborted")
		}

		// It's critical to remove/cleanup local storage
		o.cleanUp()
		o.removeStateConfig()

		// Download the current state
		err := azure.DownloadFileFromBlob(a.launchPadStorageID, o.Workspace, o.StateName+".tfstate", stateFileName)
		cobra.CheckErr(err)

		// Reset back to use local state
		console.Warning("Resetting state to local, have to re-run init without a backend/remote state")
		err = o.runLaunchpadInit(tf, true)
		cobra.CheckErr(err)
		// This is critical and stops terraform from trying to use remote state
		_ = os.Remove(o.SourcePath + "/backend.azurerm.tf")

		// Tell destroy to use local downloaded state to destroy a launchpad
		// TODO: This is a deprecated option, the solution is to switch to this
		// https://www.terraform.io/docs/language/settings/backends/local.html
		// nolint
		destroyOptions = append(destroyOptions, tfexec.State(stateFileName))
	} else {
		// Connect to launchpad, setting all the vars needed by the landingzone
		err := o.connectToLaunchPad(a.launchPadStorageID)
		cobra.CheckErr(err)
	}

	// Merge all tfvars found in config directory into -var-file options
	varOpts, err := terraform.ExpandVarDirectory(o.ConfigPath)
	cobra.CheckErr(err)
	for _, vo := range varOpts {
		// Note. spread operator would not work here, I tried ¯\_(ツ)_/¯
		destroyOptions = append(destroyOptions, vo)
	}

	console.Warning("Destroy is now running ...")
	console.StartSpinner()
	err = tf.Destroy(context.Background(), destroyOptions...)
	console.StopSpinner()
	cobra.CheckErr(err)

	// Remove files
	o.cleanUp()
	_ = os.RemoveAll(o.DataDir)

	console.Success("Destroy was successful")

	return nil
}
