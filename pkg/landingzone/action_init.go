package landingzone

import "github.com/aztfmod/rover/pkg/console"

type InitAction struct {
	TerraformAction
}

func NewInitAction() *InitAction {
	return &InitAction{
		TerraformAction{
			launchPadStorageID: "",
			ActionBase: ActionBase{
				Name:        "init",
				Description: "Perform a terraform init and no other action",
			},
		},
	}
}

func (a *InitAction) Execute(o *Options) error {
	console.Info("Carrying out Terraform init")

	a.tfexec = a.prepareTerraformCAF(o)

	if o.DryRun {
		return nil
	}

	a.runTerraformInit(o, a.tfexec)
	return nil
}