package symphony

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

func (c Config) RunAll(action landingzone.Action) {
	console.Infof("Starting CD process for all levels...\n")

	for _, level := range c.Levels {
		c.RunLevel(level, action)
	}
}

func (c Config) RunLevel(level Level, action landingzone.Action) {
	console.Infof(" - Running CD for level: %s\n", level.Name)
	for _, stack := range level.Stacks {
		c.runStack(level, &stack, action)
	}
}

// This runs the given action against the stack
// It builds a landingzone.Options struct just like landingzone.NewOptionsFromCmd() but uses the YAML as source not the cmd
func (c Config) runStack(level Level, stack *Stack, action landingzone.Action) {
	console.Infof("   - Running CD for stack: %s\n", stack.Name)

	ws := c.Workspace
	if ws == "" {
		ws = "tfstate"
	}

	cafEnv := c.Environment
	if cafEnv == "" {
		cafEnv = "sandpit"
	}

	sourcePath := c.LandingZonePath
	if sourcePath == "" {
		cobra.CheckErr("Config file is missing 'landingZonePath' key")
	}
	configPath := stack.ConfigurationPath
	if configPath == "" {
		cobra.CheckErr("Stack is missing 'configurationPath' key")
	}

	// TODO: Remove this safe guard when landingzone deploy is working
	if !level.Launchpad {
		console.Error("landingzone deployment is not implemented yet")
		return
	}

	opt := landingzone.Options{
		Level:          level.Name,
		LaunchPadMode:  level.Launchpad,
		CafEnvironment: cafEnv,
		StateName:      stack.Name,
		Workspace:      ws,
	}
	// Safely set the paths up
	opt.SetSourcePath(sourcePath)
	opt.SetConfigPath(configPath)

	// Now we can start the execution just like `landingzone run` cmd does
	opt.Execute(action)
}
