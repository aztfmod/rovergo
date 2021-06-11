package landingzone

import (
	"context"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
)

type ValidateAction struct {
	TerraformAction
}

func NewValidateAction() *ValidateAction {
	return &ValidateAction{
		TerraformAction: TerraformAction{
			launchPadStorageID: "",
			tfexec:             nil,
			ActionBase: ActionBase{
				name:        "validate",
				description: "Perform a terraform validate",
			},
		},
	}
}

func (a *ValidateAction) Execute(o *Options) error {
	console.Info("Carrying out Terraform validate")

	a.tfexec = a.prepareTerraformCAF(o)

	if o.DryRun {
		return nil
	}

	console.StartSpinner()
	out, err := a.tfexec.Validate(context.Background())
	cobra.CheckErr(err)
	console.StopSpinner()

	if !out.Valid {
		console.Errorf("Valdate returned %d warnings\n", out.WarningCount)
		console.Errorf("Valdate returned %d errors\n", out.ErrorCount)
		for _, d := range out.Diagnostics {
			console.Error("---------------")
			console.Errorf("Severity: %s\n", d.Severity)
			console.Errorf("Detail: %s\n", d.Detail)
			console.Errorf("Summary: %s\n", d.Summary)
			console.Errorf("Filename: %s\n", d.Range.Filename)
			console.Errorf("Line: %d\n", d.Range.Start.Line)
		}
		cobra.CheckErr("Validate detected issues")
	}

	console.Success("Validate was successful")
	return nil
}
