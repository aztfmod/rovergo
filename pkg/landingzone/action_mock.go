package landingzone

import "github.com/aztfmod/rover/pkg/console"

type MockAction struct {
	ActionBase
}

func (ma MockAction) Execute(o *Options) error {

	console.Infof("Environment is: %s\n", o.CafEnvironment)
	console.Infof("Config path is: %s\n", o.ConfigPath)
	console.Infof("Dry run flag is: %v\n", o.DryRun)
	console.Infof("Launchpad mode is: %v\n", o.LaunchPadMode)
	console.Infof("Level is: %s\n", o.Level)
	console.Infof("DataDir is: %s\n", o.DataDir)
	console.Infof("Source path is: %s\n", o.SourcePath)
	console.Infof("State name is: %s\n", o.StateName)
	console.Infof("State subscription is: %s\n", o.StateSubscription)
	console.Infof("Target subscription is: %s\n", o.TargetSubscription)
	console.Infof("Workspace is: %s\n", o.Workspace)
	console.Infof("Identity is: %s\n", o.Identity.DisplayName)

	return nil
}
