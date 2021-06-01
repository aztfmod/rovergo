//
// Rover - Landing zone and launchpad actions
// * Encapsulation of actions
// * Ben C, May 2021
//

package landingzone

import (
	"errors"
	"strings"
)

type Action int

var actionEnum = []string{"init", "plan", "deploy", "destroy"}

const (
	// ActionInit carries out a just init step and no real action
	ActionInit Action = iota
	// ActionPlan carries out a plan operation
	ActionPlan Action = iota
	// ActionDeploy carries out a plan AND apply operation
	ActionDeploy Action = iota
	// ActionDestroy carries out a destroy operation
	ActionDestroy Action = iota
)

func ActionFromString(actionString string) (Action, error) {
	switch strings.ToLower(actionString) {
	case ActionInit.String():
		return ActionInit, nil
	case ActionPlan.String():
		return ActionPlan, nil
	case ActionDeploy.String():
		return ActionDeploy, nil
	case ActionDestroy.String():
		return ActionDestroy, nil
	default:
		return 0, errors.New("action is not valid, must be [init | plan | deploy | destroy]")
	}
}

func (a Action) String() string {
	return actionEnum[a]
}

// func ActionToString(action Action) string {
// 	switch action {
// 	case ActionInit:
// 		return "Init"
// 	case ActionPlan:
// 		return "Plan"
// 	case ActionDeploy:
// 		return "Deploy"
// 	case ActionDestroy:
// 		return "Destroy"
// 	}
// 	return ""
// }
