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

	var err error
	a.tfexec, err = a.prepareTerraformCAF(o)
	if err != nil {
		return err
	}

	if o.DryRun {
		return nil
	}

	a.runTerraformInit(o, a.tfexec, false)
	return nil
}
