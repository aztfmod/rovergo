//
// Rover - Core functions shared by landing zone & launchpad
// * VERY WIP
// * Ben C, May 2021
//

package core

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

type Action int

const (
	// ActionInit carries out a just init step and no real action
	ActionInit Action = 1
	// ActionPlan carries out a plan operation
	ActionPlan Action = 2
	// ActionDeploy carries out a plan AND apply operation
	ActionDeploy Action = 3
	// ActionDestroy carries out a destroy operation
	ActionDestroy Action = 4
)

type Config struct {
	LaunchPadMode      bool
	ConfigPath         string
	SourcePath         string
	Level              int
	CafEnvironment     string
	StateName          string
	Workspace          string
	TargetSubscription string
	StateSubscription  string
	Account            azure.Account
}

func NewConfigFromCmd(cmd *cobra.Command) Config {
	launchPadMode := false
	if cmd.Parent().Name() == "launchpad" {
		launchPadMode = true
	}

	configPath, _ := cmd.Flags().GetString("config-path")
	sourcePath, _ := cmd.Flags().GetString("source")
	level, _ := cmd.Flags().GetInt("level")
	env, _ := cmd.Flags().GetString("environment")
	stateName, _ := cmd.Flags().GetString("statename")
	ws, _ := cmd.Flags().GetString("workspace")
	stateSub, _ := cmd.Flags().GetString("state-sub")
	targetSub, _ := cmd.Flags().GetString("target-sub")

	c := Config{
		LaunchPadMode:      launchPadMode,
		ConfigPath:         configPath,
		SourcePath:         sourcePath,
		Level:              level,
		CafEnvironment:     env,
		StateName:          stateName,
		Workspace:          ws,
		TargetSubscription: targetSub,
		StateSubscription:  stateSub,
	}

	return c
}

func RunCmd(cmd *cobra.Command, args []string) {
	actionStr, _ := cmd.Flags().GetString("action")
	action, err := ActionFromString(actionStr)
	cobra.CheckErr(err)
	console.Infof("Starting '%s' command with action '%s'\n", cmd.CommandPath(), actionStr)

	// Build config from command flags
	conf := NewConfigFromCmd(cmd)

	// Get current Azure details, subscription etc from CLI
	acct := azure.GetAccount()
	// If they weren't set as flags fall back to logged in account subscription
	if conf.StateSubscription == "" {
		conf.StateSubscription = acct.SubscriptionID
	}
	if conf.TargetSubscription == "" {
		conf.TargetSubscription = acct.SubscriptionID
	}
	conf.Account = acct

	if console.DebugEnabled {
		debugConf, _ := json.MarshalIndent(conf, "", "  ")
		console.Debug(string(debugConf))
	}

	// This should be enough for Terraform and if TargetSubscription is diff from CLI that's will work
	os.Setenv("ARM_SUBSCRIPTION_ID", conf.TargetSubscription)
	os.Setenv("TF_VAR_tfstate_subscription_id", conf.StateSubscription)

	if conf.LaunchPadMode {
		if conf.TargetSubscription != conf.StateSubscription {
			cobra.CheckErr("In launchpad mode, state-sub and target-sub must be the same Azure subscription")
		}
	}

	// Terraform init is run for all actions
	err = conf.Init()
	cobra.CheckErr(err)

	// If the action is just init, then stop here
	if action == ActionInit {
		return
	}

	err = conf.RunAction(action)
	cobra.CheckErr(err)
}

func ActionFromString(actionString string) (Action, error) {
	switch strings.ToLower(actionString) {
	case "init":
		return ActionInit, nil
	case "plan":
		return ActionPlan, nil
	case "deploy":
		return ActionDeploy, nil
	case "destroy":
		return ActionDestroy, nil
	default:
		return 0, errors.New("action is not valid, must be [init | plan | deploy | destroy]")
	}
}

// SetSharedFlags configures command flags for both landingzone and launchpad commands
func SetSharedFlags(cmd *cobra.Command) {
	cmd.PersistentFlags()
	cmd.Flags().StringP("action", "a", "init", "Action to run, one of [init | plan | deploy | destroy]")
	cmd.Flags().StringP("source", "s", "", "Path to source of landingzone (required)")
	cmd.Flags().StringP("config-path", "c", "", "Configuration vars directory (required)")
	cmd.Flags().StringP("environment", "e", "caf", "Name of CAF environment")
	cmd.Flags().StringP("workspace", "w", "tfstate", "Name of workspace")
	cmd.Flags().StringP("statename", "n", "mystate", "Base name for state and plan files")
	cmd.Flags().String("state-sub", "", "Azure subscription ID where state is held")
	cmd.Flags().String("target-sub", "", "Azure subscription ID to operate on")
	cmd.Flags().SortFlags = true

	// Level command removed from launchpad cmd as it's always zero
	if cmd.Parent().Name() != "launchpad" {
		cmd.Flags().IntP("level", "l", 1, "Level")
	}

	_ = cobra.MarkFlagRequired(cmd.Flags(), "source")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "config-path")
}

func (c Config) RunAction(action Action) error {
	console.Info("STARTING ACTION")
	return nil
}

func (c Config) Init() error {
	tfPath, err := terraform.Setup()
	cobra.CheckErr(err)
	tf, err := tfexec.NewTerraform(c.SourcePath, tfPath)
	cobra.CheckErr(err)

	if c.LaunchPadMode {
		console.Info("Running init in launchpad mode")
		err = tf.Init(context.Background(), tfexec.Upgrade(true))
		return err
	}

	console.Info("Running init for landingzone")
	// TODO: Add code to locate launchpad and storage account
	return nil
}
