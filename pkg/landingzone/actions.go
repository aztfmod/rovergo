package landingzone

import (
	"github.com/hashicorp/terraform-exec/tfexec"
)

type Action interface {
	Execute(o *Options) error
	GetName() string
	GetDescription() string
}

type ActionBase struct {
	Name        string
	Description string
}

type TerraformAction struct {
	ActionBase
	launchPadStorageID string
	tfexec             *tfexec.Terraform
}

func (ab ActionBase) GetName() string {
	return ab.Name
}

func (ab ActionBase) GetDescription() string {
	return ab.Description
}
