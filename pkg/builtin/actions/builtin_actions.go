package actions

import "github.com/aztfmod/rover/pkg/landingzone"

// ActionMap is exported so tests can use
var ActionMap = map[string]landingzone.Action{
	"init":     landingzone.NewInitAction(),
	"plan":     landingzone.NewPlanAction(),
	"apply":    landingzone.NewApplyAction(),
	"destroy":  landingzone.NewDestroyAction(),
	"validate": landingzone.NewValidateAction(),
	"fmt":      landingzone.NewFormatAction(),
	"test":     landingzone.NewTestAction(),
}
