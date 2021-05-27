//
// Rover - Landing zone run command is the core of rover
// * This carries out actions like plan, apply or destroy via terrafrom
// * TODO: !!THIS IS NOT EVEN CLOSE TO BEING COMPLETED!!
// * Ben C, May 2021
//

package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/fatih/color"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const planFileName = "rover.tfplan"

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
	lzRunCmd.Flags().StringP("action", "a", "plan", "Action to run, one of [plan | deploy | destroy]")
	lzRunCmd.Flags().StringP("source", "s", "", "Source of landingzone (required)")
	lzRunCmd.Flags().StringP("env", "e", "caf", "Name of environment")
	lzRunCmd.Flags().StringP("vars", "v", ".", "Where configuration vars are located")
	lzRunCmd.Flags().IntP("level", "l", 1, "Level")
	lzRunCmd.Flags().StringP("statename", "n", "mystate", "State and plan name")

	err := cobra.MarkFlagRequired(lzRunCmd.Flags(), "source")
	cobra.CheckErr(err)
	landingzoneCmd.AddCommand(lzRunCmd)
}

func runAction(action string, source string, varsLocation string, stateKey string, env string, level int) {
	color.Green("Running %s operation for landingzone %s", action, source)
	color.Green(" - Loading vars from: %s", varsLocation)
	color.Green(" - Level: %d", level)
	color.Green(" - State name: %s", stateKey)
	color.Green(" - Environment name: %s", env)

	tfPath, err := terraform.Setup()
	cobra.CheckErr(err)
	tf, err := tfexec.NewTerraform(source, tfPath)
	cobra.CheckErr(err)

	// TEMPORARY IMPLEMENTATION!
	// Always requires remote state config
	initOpts := []tfexec.InitOption{
		tfexec.BackendConfig(fmt.Sprintf("storage_account_name=%s", viper.GetString("state.storage-account"))),
		tfexec.BackendConfig(fmt.Sprintf("container_name=%s", viper.GetString("state.container"))),
		tfexec.BackendConfig(fmt.Sprintf("resource_group_name=%s", viper.GetString("state.resource-group"))),
		tfexec.BackendConfig(fmt.Sprintf("access_key=%s", viper.GetString("state.access-key"))),
		tfexec.BackendConfig(fmt.Sprintf("key=%s", stateKey)),
		tfexec.Reconfigure(true),
		tfexec.Upgrade(true),
		tfexec.Backend(true),
	}

	for k := range viper.GetStringMapString("state") {
		console.Debug(fmt.Sprintf("%s = %s", k, viper.GetString("state."+k)))
	}
	console.Debug(fmt.Sprintf("key = %s", stateKey))

	color.Blue("RUNNING INIT")
	err = tf.Init(context.Background(), initOpts...)
	cobra.CheckErr(err)

	switch strings.ToLower(action) {
	case "plan":
		color.Blue("RUNNING PLAN")
		result, err := tf.Plan(context.Background(), tfexec.Out(planFileName))
		color.Blue("PLAN RESULT WAS %v", result)
		cobra.CheckErr(err)
	case "apply":
		color.Blue("RUNNING APPLY")
		err := tf.Apply(context.Background(), tfexec.DirOrPlan(planFileName))
		cobra.CheckErr(err)
	default:
		cobra.CheckErr(color.RedString("provided action '%s' is invalid", action))
	}

}
