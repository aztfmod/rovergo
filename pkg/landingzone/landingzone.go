package landingzones

import "github.com/aztfmod/rover/pkg/console"

type Action int

const (
	// ActionPlan carries out a plan operation
	ActionPlan Action = 1
	// ActionDeploy carries out a plan AND apply operation
	ActionDeploy Action = 2
	// ActionDestroy carries out a destroy operation
	ActionDestroy Action = 3
)

type Config struct {
	LaunchPadMode  bool
	ConfigPath     string
	SourcePath     string
	Level          int
	CafEnvironment string
	StateName      string
	Workspace      string
	Subscription   string
}

func (c Config) Run(action Action) error {
	console.Info("STARTING ACTION")
	return nil
}

func (c Config) Init() error {
	if c.LaunchPadMode {
		console.Debug("INIT LAUNCHPAD")
		return nil
	}

	console.Debug("INIT LANDINGZONE")
	return nil
}
