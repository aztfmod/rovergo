package symphony

import (
	"fmt"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

func BuildOptions(cmd *cobra.Command) []landingzone.Options {
	launchPadMode, _ := cmd.Flags().GetBool("launchpad")
	configFile, _ := cmd.Flags().GetString("config-file")
	sourcePath, _ := cmd.Flags().GetString("source")
	levelName, _ := cmd.Flags().GetString("level")
	env, _ := cmd.Flags().GetString("environment")
	stateName, _ := cmd.Flags().GetString("statename")
	ws, _ := cmd.Flags().GetString("workspace")
	stateSub, _ := cmd.Flags().GetString("state-sub")
	targetSub, _ := cmd.Flags().GetString("target-sub")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if launchPadMode || env != "" || ws != "" || stateName != "" || stateSub != "" || targetSub != "" || sourcePath != "" {
		cobra.CheckErr("Do not supply any options other than level when using a config file")
	}

	conf, err := NewSymphonyConfig(configFile)
	cobra.CheckErr(err)

	if levelName == "" {
		console.Info("Rover will operate on ALL levels...")
		isDestroy := cmd.Name() == "destroy"
		optionsList := conf.parseAllLevels(isDestroy)
		// We munge some options here rather than passing it through all the parser functions
		for i := range optionsList {
			optionsList[i].DryRun = dryRun
		}
		return optionsList
	}

	var level *Level
	// Try to locate level in config matching the level flag passed in
	for _, confLevel := range conf.Content.Levels {
		if confLevel.Name == levelName {
			level = &confLevel
			break
		}
	}

	// nolint
	if level == nil {
		cobra.CheckErr(fmt.Sprintf("level '%s' not found in symphony config file", levelName))
	}

	// nolint
	console.Infof("Rover will operate on level '%s'...\n", level.Name)
	// nolint
	optionsList := conf.parseLevel(*level)

	// We munge some options here rather than passing it through all the parser functions
	for i := range optionsList {
		optionsList[i].DryRun = dryRun
	}

	return optionsList
}

func SetDry(o *landingzone.Options) {
	o.DryRun = true
}
