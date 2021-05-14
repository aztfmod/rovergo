package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var lzRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an action for landingzones",

	Run: func(cmd *cobra.Command, args []string) {

		action, _ := cmd.Flags().GetString("action")
		source, _ := cmd.Flags().GetString("source")
		stateKey, _ := cmd.Flags().GetString("statename")
		env, _ := cmd.Flags().GetString("env")
		vars, _ := cmd.Flags().GetString("vars")
		level, _ := cmd.Flags().GetInt("level")

		runAction(action, source, vars, stateKey, env, level)
	},
}

func init() {
	lzRunCmd.PersistentFlags()
	lzRunCmd.Flags().StringP("action", "a", "plan", "Action to run, one of [plan | apply | destroy]")
	lzRunCmd.Flags().StringP("source", "s", "", "Source of landingzone (required)")
	lzRunCmd.Flags().StringP("env", "e", "caf", "Name of environment")
	lzRunCmd.Flags().StringP("vars", "v", ".", "Where configuration vars are located")
	lzRunCmd.Flags().IntP("level", "l", 1, "Level")
	lzRunCmd.Flags().StringP("statename", "n", "mystate", "State and plan name")

	cobra.MarkFlagRequired(lzRunCmd.Flags(), "source")
	landingzoneCmd.AddCommand(lzRunCmd)
}

func runAction(action string, source string, varsLocation string, stateKey string, env string, level int) {
	color.Green("Running %s operation for landingzone %s", action, source)
	color.Green(" - Loading vars from: %s", varsLocation)
	color.Green(" - Level: %d", level)
	color.Green(" - State name: %s", stateKey)
	color.Green(" - Environment name: %s", env)
}
