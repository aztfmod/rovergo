package landingzone

import (
	"context"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type FormatAction struct {
	TerraformAction
}

func NewFormatAction() *FormatAction {
	return &FormatAction{
		TerraformAction: TerraformAction{
			launchPadStorageID: "",
			ActionBase: ActionBase{
				Name:        "fmt",
				Type:        BuiltinCommand,
				Description: "Perform a terraform format",
			},
		},
	}
}

func (a *FormatAction) Execute(o *Options) error {
	tf, err := a.prepareTerraformCAF(o)
	if err != nil {
		return err
	}

	if o.DryRun {
		return nil
	}

	fo := []tfexec.FormatOption{
		tfexec.Dir(o.SourcePath),
		tfexec.Recursive(true),
	}

	outcome, filesToFix, err := tf.FormatCheck(context.Background(), fo...)
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
