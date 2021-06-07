//
// Rover - Top level cd (continuous deployment) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"github.com/spf13/cobra"
)

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Manage CD operations.",
	Long:  `Manage CD operations.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	levelFlag, _ := cmd.Flags().GetString("level")
	// 	configPath, _ := cmd.Flags().GetString("symphony-config")
	// 	actionStr, _ := cmd.Flags().GetString("action")
	// 	action, err := landingzone.NewAction(actionStr)
	// 	cobra.CheckErr(err)

	// 	conf, err := symphony.NewSymphonyConfig(configPath)
	// 	cobra.CheckErr(err)

	// 	if levelFlag == "" {
	// 		conf.RunAll(action)
	// 		return
	// 	}

	// 	var level *symphony.Level
	// 	// Try to locate level in config matching the level flag passed in
	// 	for _, confLevel := range conf.Content.Levels {
	// 		if confLevel.Name == levelFlag {
	// 			level = &confLevel
	// 			break
	// 		}
	// 	}

	// 	//nolint
	// 	if level == nil {
	// 		cobra.CheckErr(fmt.Sprintf("level '%s' not found in symphony config file", levelFlag))
	// 	}
	// 	//nolint
	// 	conf.RunLevel(*level, action)
	// },
}

func init() {
	rootCmd.AddCommand(cdCmd)
	// cdCmd.Flags().StringP("symphony-config", "c", "", "Path to symphony config file")
	// cdCmd.Flags().StringP("level", "l", "", "Level to operate on, if omitted all levels will be processed")
	// landingzone.AddActionFlag(cdCmd)

	// _ = cobra.MarkFlagRequired(cdCmd.Flags(), "symphony-config")
}
