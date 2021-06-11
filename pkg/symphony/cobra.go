package symphony

import (
	"fmt"

	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

func RunFromConfig(cmd *cobra.Command, action landingzone.Action) {
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
		conf.runAll(action, dryRun)
		return
	}

	var level *Level
	// Try to locate level in config matching the level flag passed in
	for _, confLevel := range conf.Content.Levels {
		if confLevel.Name == levelName {
			level = &confLevel
			break
		}
	}

	//nolint
	if level == nil {
		cobra.CheckErr(fmt.Sprintf("level '%s' not found in symphony config file", levelName))
	}

	//nolint
	conf.RunLevel(*level, action, dryRun)
}
