package symphony

import (
	"path/filepath"

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
	console.Infof(" - Running CD for level %d\n", level.Number)
	for _, stack := range level.Stacks {
		c.runStack(level, &stack, action)
	}
}

func (c Config) runStack(level Level, stack *Stack, action landingzone.Action) {
	console.Infof("   - Running CD for stack %s\n", stack.Name)

	ws := c.Workspace
	if ws == "" {
		ws = "tfstate"
	}

	cafEnv := c.Environment
	if cafEnv == "" {
		cafEnv = "sandpit"
	}

	lzPath := c.LandingZonePath
	if lzPath == "" {
		cobra.CheckErr("Config file is missing 'landingZonePath' setting")
	}

	// Convert to absolute paths as a precaution
	lzPath, err := filepath.Abs(lzPath)
	cobra.CheckErr(err)
	lzPath = landingzone.SetSourceDir(lzPath, level.Launchpad)
	configPath, err := filepath.Abs(stack.ConfigurationPath)
	cobra.CheckErr(err)

	// TODO: Remove this safe guard when landingzone deploy is working
	if !level.Launchpad {
		console.Error("landingzone deployment is not implemented yet")
		return
	}

	opt := landingzone.Options{
		Level:          level.Number,
		ConfigPath:     configPath,
		SourcePath:     lzPath,
		LaunchPadMode:  level.Launchpad,
		CafEnvironment: cafEnv,
		StateName:      stack.Name,
		Workspace:      ws,
	}

	landingzone.ExecuteRun(opt, action)
}
