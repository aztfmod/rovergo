//
// Rover - Core functions shared by landing zone & launchpad
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/briandowns/spinner"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

// RunCmd is the shared cobra run function for `landingzone run` and `launchpad run`
func RunCmd(cmd *cobra.Command, args []string) {
	actionStr, _ := cmd.Flags().GetString("action")
	action, err := ActionFromString(actionStr)
	cobra.CheckErr(err)
	console.Infof("Starting '%s' command with action '%s'\n", cmd.CommandPath(), actionStr)

	// Build config from command flags
	opt := NewOptionsFromCmd(cmd)

	// Get current Azure details, subscription etc from CLI
	acct := azure.GetSubscription()
	ident := azure.GetIdentity()
	// If they weren't set as flags fall back to logged in account subscription
	if opt.StateSubscription == "" {
		opt.StateSubscription = acct.ID
	}
	if opt.TargetSubscription == "" {
		opt.TargetSubscription = acct.ID
	}
	opt.Subscription = acct
	opt.Identity = ident

	if opt.LaunchPadMode {
		if opt.TargetSubscription != opt.StateSubscription {
			cobra.CheckErr("In launchpad mode, state-sub and target-sub must be the same Azure subscription")
		}
	}

	// Remove old files, reset backend etc
	opt.cleanUp()

	// All the env vars & other setup needed for CAF and get handle on Terraform
	tf := opt.initializeCAF(cmd.Root().Version)

	// Find state storage account for this environment and level
	storageID, err := azure.FindStorageAccount(opt.LevelString(), opt.CafEnvironment, opt.StateSubscription)
	if err != nil {
		if opt.LaunchPadMode {
			console.Warning("No state storage account found, but running in launchpad mode, it will be created")
		} else {
			console.Errorf("No state storage account found for environment '%s' and level %d, please deploy a launchpad first!\n", opt.CafEnvironment, opt.Level)
			cobra.CheckErr("Can't deploy a landing zone without a launchpad")
		}
	} else {
		console.Infof("Located state storage account %s\n", storageID)
	}

	// Run init in correct mode
	if opt.LaunchPadMode && storageID == "" {
		err = opt.runLaunchpadInit(tf)
	} else {
		err = opt.runRemoteInit(tf, storageID)
	}
	cobra.CheckErr(err)

	// If the action is just init, then stop here
	if action == ActionInit {
		return
	}

	err = opt.runAction(tf, action)
	cobra.CheckErr(err)

	// Special case for post launchpad deployment
	newStorageID, err := azure.FindStorageAccount(opt.LevelString(), opt.CafEnvironment, opt.StateSubscription)
	if opt.LaunchPadMode && storageID != newStorageID {
		console.Info("Detected the launchpad infrastructure has been deployed or updated")
		cobra.CheckErr(err)
		stateFileName := opt.OutPath + "/" + opt.StateName + ".tfstate"
		azure.UploadFileToBlob(newStorageID, opt.Workspace, opt.StateName+".tfstate", stateFileName)
		console.Info("Uploading state from launchpad process to Azure storage")
		os.Remove(stateFileName)
	}

	console.Info("Rover completed")
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

// Carry out the plan/deploy/destroy action
func (o Options) runAction(tf *tfexec.Terraform, action Action) error {
	console.Infof("Starting '%s' action, this could take some time...\n", action.String())
	spinner := spinner.New(spinner.CharSets[37], 100*time.Millisecond)

	// Plan is run for both plan and deploy actions
	if action == ActionPlan || action == ActionDeploy {
		console.Info("Carrying out the Terraform plan phase")

		// Build plan options starting with tfplan output
		planFile := fmt.Sprintf("%s/%s.tfplan", o.OutPath, o.StateName)
		planOptions := []tfexec.PlanOption{
			tfexec.Out(planFile),
			tfexec.Refresh(true),
		}

		// Then merge all tfvars found in config directory into -var-file options
		varOpts, err := terraform.ExpandVarDirectory(o.ConfigPath)
		cobra.CheckErr(err)
		planOptions = append(planOptions, varOpts...)

		// Now actually invoke Terraform plan
		spinner.Start()
		changes, err := tf.Plan(context.Background(), planOptions...)
		spinner.Stop()
		if err != nil {
			return err
		}
		if changes {
			console.Success("Plan contains infrastructure updates")
		} else {
			console.Success("Plan detected no changes")
			console.Success("Skipping the apply phase")
			return nil
		}
	}

	if action == ActionDeploy {
		console.Info("Carrying out the Terraform apply phase")

		planFile := fmt.Sprintf("%s/%s.tfplan", o.OutPath, o.StateName)
		stateFile := fmt.Sprintf("%s/%s.tfstate", o.OutPath, o.StateName)

		// Build apply options, with plan file and state out
		applyOptions := []tfexec.ApplyOption{
			tfexec.DirOrPlan(planFile),
			tfexec.StateOut(stateFile),
		}

		// Now actually invoke Terraform apply
		spinner.Start()
		err := tf.Apply(context.Background(), applyOptions...)
		spinner.Stop()
		if err != nil {
			return err
		}

		console.Success("Apply was successful")
	}

	return nil
}

// Carry out Terraform init operation in launchpad mode has no backend state
func (o Options) runLaunchpadInit(tf *tfexec.Terraform) error {
	console.Info("Running init for launchpad")
	spinner := spinner.New(spinner.CharSets[37], 100*time.Millisecond)

	spinner.Start()
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	spinner.Stop()
	return err
}

// Carry out Terraform init operation with remote state backend
func (o Options) runRemoteInit(tf *tfexec.Terraform, storageID string) error {
	console.Info("Running init with remote state")
	spinner := spinner.New(spinner.CharSets[37], 100*time.Millisecond)

	// IMPORTANT: This enables remote state in the source terraform dir
	o.enableAzureBackend()

	subID, resGrp, accountName := azure.ParseResourceID(storageID)
	accessKey := azure.GetAccountKey(subID, accountName, resGrp)

	initOptions := []tfexec.InitOption{
		tfexec.BackendConfig(fmt.Sprintf("storage_account_name=%s", accountName)),
		tfexec.BackendConfig(fmt.Sprintf("container_name=%s", o.Workspace)),
		tfexec.BackendConfig(fmt.Sprintf("resource_group_name=%s", resGrp)),
		tfexec.BackendConfig(fmt.Sprintf("access_key=%s", accessKey)),
		tfexec.BackendConfig(fmt.Sprintf("key=%s", o.StateName+".tfstate")),
		tfexec.Reconfigure(true),
		tfexec.Upgrade(true),
		tfexec.Backend(true),
	}

	spinner.Start()
	err := tf.Init(context.Background(), initOptions...)
	cobra.CheckErr(err)
	spinner.Stop()
	return err
}

// All env vars and other steps need before running Terraform with CAF landingzones
func (o *Options) initializeCAF(roverVersion string) *tfexec.Terraform {
	tfPath, err := terraform.Setup()
	cobra.CheckErr(err)

	os.Setenv("ARM_SUBSCRIPTION_ID", o.TargetSubscription)
	os.Setenv("ARM_TENANT_ID", o.Subscription.TenantID)
	os.Setenv("TF_VAR_tfstate_subscription_id", o.StateSubscription)
	os.Setenv("TF_VAR_tf_name", fmt.Sprintf("%s.tfstate", o.StateName))
	os.Setenv("TF_VAR_tf_plan", fmt.Sprintf("%s.tfplan", o.StateName))
	os.Setenv("TF_VAR_workspace", o.Workspace)
	os.Setenv("TF_VAR_level", o.LevelString())
	os.Setenv("TF_VAR_environment", o.CafEnvironment)
	os.Setenv("TF_VAR_rover_version", roverVersion)
	os.Setenv("TF_VAR_tenant_id", o.Subscription.TenantID)
	os.Setenv("TF_VAR_user_type", o.Identity.ObjectType)
	os.Setenv("TF_VAR_logged_user_objectId", o.Identity.ObjectID)

	// TODO: Removed for now pending further investigation
	// envName := o.Account.EnvironmentName
	// // For some reason the name returned from the CLI for Azure public is not valid!
	// if envName == "AzureCloud" {
	// 	envName = "AzurePublicCloud"
	// }
	// os.Setenv("AZURE_ENVIRONMENT", envName)
	// os.Setenv("ARM_ENVIRONMENT", azure.CloudNameToTerraform(envName))

	// Default the TF_DATA_DIR to user's home dir
	dataDir := os.Getenv("TF_DATA_DIR")
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		os.Setenv("TF_DATA_DIR", home)
	}

	// Create local state/plan folder, rover puts this in a opinionated place, for reasons I don't understand
	localStatePath := fmt.Sprintf("%s/tfstates/%s/%s", os.Getenv("TF_DATA_DIR"), o.LevelString(), o.Workspace)
	err = os.MkdirAll(localStatePath, os.ModePerm)
	cobra.CheckErr(err)
	o.OutPath = localStatePath

	// Create new TF exec with the working dir set to source
	tf, err := tfexec.NewTerraform(o.SourcePath, tfPath)
	cobra.CheckErr(err)

	// The debugging done here
	if console.DebugEnabled {
		// By default no output from Terraform is seen
		// Also noote TF_LOG env var is ignored by tfexec
		// tf.SetStdout(os.Stdout)
		// tf.SetStderr(os.Stderr)

		// This gives us some info level logs we can send to stdout
		tf.SetLogger(console.Printfer{})

		console.Debug("##### DEBUG OPTIONS:")
		debugConf, _ := json.MarshalIndent(o, "", "  ")
		console.Debug(string(debugConf))

		console.Debug("##### DEBUG ENV VARS:")
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "ARM_") || strings.HasPrefix(env, "AZURE_") || strings.HasPrefix(env, "TF_") {
				console.Debug(env)
			}
		}
	}
	return tf
}

// Remove files to ensure a clean run, state and plan files are recreated
func (o Options) cleanUp() {
	_ = os.Remove(o.SourcePath + "/backend.azurerm.tf")
	_ = os.Remove(o.OutPath + "/" + o.StateName + ".tfstate")
	_ = os.Remove(o.OutPath + "/" + o.StateName + ".tfplan")
	_ = os.Remove(os.Getenv("TF_DATA_DIR") + "/terraform.tfstate")
}

// By copying this file we enable teh azurerm backend and therefore remote state
func (o Options) enableAzureBackend() {
	console.Info("Enabling backend state with backend.azurerm.tf file")
	err := utils.CopyFile(o.SourcePath+"/backend.azurerm", o.SourcePath+"/backend.azurerm.tf")
	cobra.CheckErr(err)
}
