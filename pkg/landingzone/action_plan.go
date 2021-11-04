package landingzone

import (
	"context"
	"fmt"
	"path"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type PlanAction struct {
	TerraformAction
	hasChanges bool
}

func NewPlanAction() *PlanAction {
	return &PlanAction{
		hasChanges: false,
		TerraformAction: TerraformAction{
			launchPadStorageID: "",
			ActionBase: ActionBase{
				Name:        "plan",
				Type:        BuiltinCommand,
				Description: "Perform a terraform plan",
			},
		},
	}
}

func (a *PlanAction) Execute(o *Options) error {
	tf, err := a.prepareTerraformCAF(o)
	if err != nil {
		return err
	}

	if o.DryRun {
		return nil
	}

	// Connect to launchpad, setting all the vars needed by the landingzone
	if !o.LaunchPadMode {
		err := o.connectToLaunchPad(a.launchPadStorageID)
		cobra.CheckErr(err)
	}

	// Build plan options starting with tfplan output
	planFile := path.Join(o.DataDir, fmt.Sprintf("%s.tfplan", o.StateName))
	planOptions := []tfexec.PlanOption{
		tfexec.Out(planFile),
		tfexec.Refresh(true),
		tfexec.Parallelism(terraformParallelism),
	}

	// Then merge all tfvars found in config directory into -var-file options
	varOpts, err := terraform.ExpandVarDirectory(o.ConfigPath)
	cobra.CheckErr(err)
	for _, vo := range varOpts {
		// Note. spread operator would not work here, I tried ¯\_(ツ)_/¯
		planOptions = append(planOptions, vo)
	}

	console.StartSpinner()
	a.hasChanges, err = tf.Plan(context.Background(), planOptions...)
	console.StopSpinner()
	cobra.CheckErr(err)
	if a.hasChanges {
		console.Successf("Plan %s contains infrastructure updates\n", planFile)
	} else {
		console.Successf("Plan %s detected no changes\n", planFile)
	}
	return nil
}
