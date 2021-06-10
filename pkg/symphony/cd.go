package symphony

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

func (c Config) RunAll(action landingzone.Action) {
	console.Infof("Starting CD process for all levels...\n")

	// Special case, handle destroy all levels in REVERSE order
	if action == landingzone.ActionDestroy {
		console.Warningf("Destroying ALL levels (in reverse order), I hope you know what you are doing...\n")
		for l := len(c.Content.Levels) - 1; l >= 0; l-- {
			c.RunLevel(c.Content.Levels[l], action)
		}
		return
	}

	for _, level := range c.Content.Levels {
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

	ws := c.Content.Workspace
	if ws == "" {
		ws = "tfstate"
	}

	cafEnv := c.Content.Environment
	if cafEnv == "" {
		cafEnv = "sandpit"
	}

	sourcePath := stack.LandingZonePath
	if sourcePath == "" {
		cobra.CheckErr("Stack is missing 'landingZonePath' key")
	}
	configPath := stack.ConfigurationPath
	if configPath == "" {
		cobra.CheckErr("Stack is missing 'configurationPath' key")
	}

	// TODO: Remove this safe guard when landingzone deploy is working
	// if !level.Launchpad {
	// 	console.Error("landingzone deployment is not implemented yet")
	// 	return
	// }

	stateName := stack.TfState
	// IMPORTANT: We use the stack name as the default name if tfState key is not supplied
	if stateName == "" {
		stateName = stack.Name
	}

	opt := landingzone.Options{
		Level:          level.Name,
		LaunchPadMode:  level.Launchpad,
		CafEnvironment: cafEnv,
		StateName:      stateName,
		Workspace:      ws,
	}

	// Safely set the paths up
	opt.SetSourcePath(sourcePath)
	opt.SetConfigPath(configPath)

	// Now we can start the execution just like `landingzone run` cmd does
	opt.Execute(action)
	console.Successf("Finished execution on stack '%s'\n", stack.Name)
}
