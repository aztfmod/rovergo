package landingzone

import (
	"github.com/hashicorp/terraform-exec/tfexec"
)

type Action interface {
	Execute(o *Options) error
	Name() string
	Description() string
}

type ActionBase struct {
	name        string
	description string
}

type TerraformAction struct {
	ActionBase
	launchPadStorageID string
	tfexec             *tfexec.Terraform
}

func (ab ActionBase) Name() string {
	return ab.name
}

func (ab ActionBase) Description() string {
	return ab.description
}
