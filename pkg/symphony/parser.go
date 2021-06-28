package symphony

import (
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

// This parses ALL levels returning a slice of Options structs one for each level and stack
func (c Config) parseAllLevels(isDestroy bool) []landingzone.Options {
	optionsList := []landingzone.Options{}

	// Special case, handle destroy all levels in REVERSE order
	if isDestroy {
		console.Warningf("Destroying ALL levels (in reverse order), I hope you know what you are doing...\n")
		for l := len(c.Content.Levels) - 1; l >= 0; l-- {
			optionsList = append(optionsList, c.parseLevel(c.Content.Levels[l])...)
		}
		return optionsList
	}

	// Normal order
	for _, level := range c.Content.Levels {
		optionsList = append(optionsList, c.parseLevel(level)...)
	}

	return optionsList
}

// This parses a level returning a slice of Options structs one for each stack
// All stacks are parsed within the level
func (c Config) parseLevel(level Level) []landingzone.Options {
	console.Infof(" - Parsing level: %s\n", level.Name)
	optionsList := []landingzone.Options{}
	for _, stack := range level.Stacks {
		optionsList = append(optionsList, c.parseStack(level, &stack))
	}
	return optionsList
}

// This parses a stack returning an Options struct that can be used to execute an Action on that stack
func (c Config) parseStack(level Level, stack *Stack) landingzone.Options {
	console.Infof("   - Parsing stack: %s\n", stack.Name)

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
	}

	// Safely set the paths up
	opt.SetSourcePath(sourcePath)
	opt.SetConfigPath(configPath)
	err := opt.SetDataDir()
	cobra.CheckErr(err)

	return opt
}
