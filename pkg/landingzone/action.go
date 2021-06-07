//
// Rover - All landingzone terraform actions
// * Encapsulation of actions
// * Ben C, May 2021
//

package landingzone

import (
	"errors"
	"strings"
)

type Action int

// ActionEnum is the list of all action strings, note order MUST match those defined as consts below
var ActionEnum = []string{"init", "plan", "apply", "run", "destroy", "test", "fmt", "validate"}

// Used when building commands
var descriptionEnum = []string{
	"Perform a terraform init and no other action",
	"Perform a terraform plan",
	"Perform a terraform plan & apply",
	"Perform a terraform plan, apply & test",
	"Perform a terraform destroy",
	"Run a test using terratest",
	"Perform a terraform fmt check",
	"Perform a terraform validate",
}

const (
	// ActionInit carries out a init operation and exits
	ActionInit Action = iota
	// ActionPlan carries out a plan operation
	ActionPlan Action = iota
	// ActionApply carries out a plan AND apply
	ActionApply Action = iota
	// ActionRun carries out a plan, apply + test
	ActionRun Action = iota
	// ActionDestroy carries out a destroy operation
	ActionDestroy Action = iota
	// ActionTest carries out a test operation
	ActionTest Action = iota
	// ActionFormat carries out a fmt operation
	ActionFormat Action = iota
	// ActionValidate carries out a vaildate operation
	ActionValidate Action = iota
)

// NewAction returns an Action type from a string
func NewAction(actionString string) (Action, error) {
	switch strings.ToLower(actionString) {
	case ActionInit.Name():
		return ActionInit, nil
	case ActionPlan.Name():
		return ActionPlan, nil
	case ActionApply.Name():
		return ActionApply, nil
	case ActionRun.Name():
		return ActionRun, nil
	case ActionDestroy.Name():
		return ActionDestroy, nil
	case ActionTest.Name():
		return ActionTest, nil
	case ActionFormat.Name():
		return ActionFormat, nil
	case ActionValidate.Name():
		return ActionValidate, nil
	default:
		return 0, errors.New("action is not valid")
	}
}

func (a Action) Name() string {
	return ActionEnum[a]
}

func (a Action) Description() string {
	return descriptionEnum[a]
}
