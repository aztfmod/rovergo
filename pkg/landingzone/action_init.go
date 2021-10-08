package landingzone

type InitAction struct {
	TerraformAction
}

func NewInitAction() *InitAction {
	return &InitAction{
		TerraformAction{
			launchPadStorageID: "",
			ActionBase: ActionBase{
				Name:        "init",
				Type:        BuiltinCommand,
				Description: "Perform a terraform init and no other action",
			},
		},
	}
}

func (a *InitAction) Execute(o *Options) error {
	tf, err := a.prepareTerraformCAF(o)
	if err != nil {
		return err
	}

	if o.DryRun {
		return nil
	}

	a.runTerraformInit(o, tf, false)
	return nil
}
