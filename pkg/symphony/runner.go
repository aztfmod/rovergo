package symphony

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

func (c Config) runAll(action landingzone.ActionI, dryRun bool) {
	console.Infof("Starting process for all levels...\n")

	// Special case, handle destroy all levels in REVERSE order
	if action.Name() == "destroy" {
		console.Warningf("Destroying ALL levels (in reverse order), I hope you know what you are doing...\n")
		for l := len(c.Content.Levels) - 1; l >= 0; l-- {
			c.RunLevel(c.Content.Levels[l], action, dryRun)
		}
		return
	}

	for _, level := range c.Content.Levels {
		c.RunLevel(level, action, dryRun)
	}
}

func (c Config) RunLevel(level Level, action landingzone.ActionI, dryRun bool) {
	console.Infof(" - Running CD for level: %s\n", level.Name)
	for _, stack := range level.Stacks {
		c.runStack(level, &stack, action, dryRun)
	}
}

// This runs the given action against the stack
// It builds a landingzone.Options struct just like landingzone.NewOptionsFromCmd() but uses the YAML as source not the cmd
func (c Config) runStack(level Level, stack *Stack, action landingzone.ActionI, dryRun bool) {
	console.Infof("   - Running stack: %s\n", stack.Name)

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

	stateName := stack.TfState
	// NOTE! We use the stack name as the default name if tfState key is not supplied
	if stateName == "" {
		stateName = stack.Name
	}

	opt := landingzone.Options{
		Level:          level.Name,
		LaunchPadMode:  level.Launchpad,
		CafEnvironment: cafEnv,
		StateName:      stateName,
		Workspace:      ws,
		DryRun:         dryRun,
	}

	// Safely set the paths up
	opt.SetSourcePath(sourcePath)
	opt.SetConfigPath(configPath)

	// Now we can start the execution just like the CLI would
	// Note errors are ignored they currently handled internally by the action
	_ = action.Execute(&opt)

	console.Successf("Finished execution on stack '%s'\n", stack.Name)
}
