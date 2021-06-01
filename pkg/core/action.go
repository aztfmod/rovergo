//
// Rover - Core actions
// * Encapsulation of actions
// * Ben C, May 2021
//

package core

import (
	"errors"
	"strings"
)

type Action int

const (
	// ActionInit carries out a just init step and no real action
	ActionInit Action = 1
	// ActionPlan carries out a plan operation
	ActionPlan Action = 2
	// ActionDeploy carries out a plan AND apply operation
	ActionDeploy Action = 3
	// ActionDestroy carries out a destroy operation
	ActionDestroy Action = 4
)

func ActionFromString(actionString string) (Action, error) {
	switch strings.ToLower(actionString) {
	case "init":
		return ActionInit, nil
	case "plan":
		return ActionPlan, nil
	case "deploy":
		return ActionDeploy, nil
	case "destroy":
		return ActionDestroy, nil
	default:
		return 0, errors.New("action is not valid, must be [init | plan | deploy | destroy]")
	}
}

func ActionToString(action Action) string {
	switch action {
	case ActionInit:
		return "Init"
	case ActionPlan:
		return "Plan"
	case ActionDeploy:
		return "Deploy"
	case ActionDestroy:
		return "Destroy"
	}
	return ""
}
