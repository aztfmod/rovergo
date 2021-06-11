package landingzone

import (
	"context"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type FormatAction struct {
	CAFAction
}

func NewFormatAction() *FormatAction {
	return &FormatAction{
		CAFAction: CAFAction{
			launchPadStorageID: "",
			tfexec:             nil,
			ActionBase: ActionBase{
				name:        "fmt",
				description: "Perform a terraform format",
			},
		},
	}
}

func (a *FormatAction) Execute(o *Options) error {
	console.Info("Carrying out Terraform format")

	a.tfexec = a.prepareTerraformCAF(o)

	if o.DryRun {
		return nil
	}

	fo := []tfexec.FormatOption{
		tfexec.Dir(o.SourcePath),
		tfexec.Recursive(true),
	}

	outcome, filesToFix, err := a.tfexec.FormatCheck(context.Background(), fo...)
	cobra.CheckErr(err)

	// TODO: return something (exit code?) so that pipeline can react appropriately
	if outcome {
		console.Success("No formatting is necessary.")
	} else {
		console.Error("The following file(s) require formatting:")
		for _, filename := range filesToFix {
			console.Errorf("  %s\n", filename)
		}
		cobra.CheckErr("Format detected issues")
	}

	console.Success("Format completed")
	return nil
}
