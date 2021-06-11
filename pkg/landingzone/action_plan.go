package landingzone

import (
	"context"
	"fmt"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type PlanAction struct {
	CAFAction
	hasChanges bool
}

func NewPlanAction() *PlanAction {
	return &PlanAction{
		hasChanges: false,
		CAFAction: CAFAction{
			launchPadStorageID: "",
			tfexec:             nil,
			ActionBase: ActionBase{
				name:        "plan",
				description: "Perform a terraform plan",
			},
		},
	}
}

func (a *PlanAction) Execute(o *Options) error {
	a.tfexec = a.prepareTerraformCAF(o)

	if o.DryRun {
		return nil
	}

	a.runTerraformInit(o, a.tfexec)

	console.Info("Carrying out Terraform plan")

	// Connect to launchpad, setting all the vars needed by the landingzone
	if !o.LaunchPadMode {
		err := o.connectToLaunchPad(a.launchPadStorageID)
		cobra.CheckErr(err)
	}

	// Build plan options starting with tfplan output
	planFile := fmt.Sprintf("%s/%s.tfplan", o.OutPath, o.StateName)
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
	a.hasChanges, err = a.tfexec.Plan(context.Background(), planOptions...)
	console.StopSpinner()
	cobra.CheckErr(err)
	if a.hasChanges {
		console.Success("Plan contains infrastructure updates")
	} else {
		console.Success("Plan detected no changes")
	}
	return nil
}
