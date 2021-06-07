//
// Rover - Core execution of landingzone operations and actions
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/aztfmod/rover/pkg/version"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

const terraformParallelism = 30

// Execute is entry point for `landingzone run`, `launchpad run` and `cd` operations
// This executes an action against a set of config options
func (o *Options) Execute(action Action) {
	// Get current Azure details, subscription etc from CLI
	acct := azure.GetSubscription()
	ident := azure.GetIdentity()
	// If they weren't set as flags fall back to logged in account subscription
	if o.StateSubscription == "" {
		o.StateSubscription = acct.ID
	}
	if o.TargetSubscription == "" {
		o.TargetSubscription = acct.ID
	}
	o.Subscription = acct
	o.Identity = ident

	if o.LaunchPadMode {
		if o.TargetSubscription != o.StateSubscription {
			cobra.CheckErr("In launchpad mode, state-sub and target-sub must be the same Azure subscription")
		}
	}

	// Remove old files, reset backend etc
	o.cleanUp()

	// All the env vars & other setup needed for CAF and get handle on Terraform
	tf := o.initializeCAF()

	// Find state storage account for this environment and level
	existingStorageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	if err != nil {
		if o.LaunchPadMode {
			console.Warning("No state storage account found, but running in launchpad mode, it will be created")
		} else {
			console.Errorf("No state storage account found for environment '%s' and level %d, please deploy a launchpad first!\n", o.CafEnvironment, o.Level)
			cobra.CheckErr("Can't deploy a landing zone without a launchpad")
		}
	} else {
		console.Infof("Located state storage account %s\n", existingStorageID)
	}

	// TODO: PUT COMMMANDS HERE THAT DONT NEED INIT AND EXIT EARLY
	if action == ActionFmt {
		console.Info("Carrying out the Terraform fmt command")

		fo := []tfexec.FormatOption{
			tfexec.Dir(o.SourcePath),
			tfexec.Recursive(true),
		}
		outcome, filesToFix, err := tf.FormatCheck(context.Background(), fo...)
		cobra.CheckErr(err)

		// TODO: return something (exit code?) so that pipeline can react appropriately
		if outcome {
			console.Success("No formatting is necessary.")
		} else {
			console.Info("The following file(s) require formatting:\n")
			for _, filename := range filesToFix {
				console.Infof("%s\n", filename)
			}
		}
		return
	}

	// Run init in correct mode
	if o.LaunchPadMode && existingStorageID == "" {
		err = o.runLaunchpadInit(tf)
	} else {
		err = o.runRemoteInit(tf, existingStorageID)
	}
	cobra.CheckErr(err)

	// If the action is just init, then stop here and don't proceed
	if action == ActionInit {
		console.Success("Rover completed, only init was run and no infrastructure changes were planned or applied")
		return
	}

	console.Infof("Starting '%s' action, this could take some time...\n", action.String())

	// Plan is run for both plan and deploy actions
	if action == ActionPlan || action == ActionDeploy {
		console.Info("Carrying out the Terraform plan phase")

		// Build plan options starting with tfplan output
		planFile := fmt.Sprintf("%s/%s.tfplan", o.OutPath, o.StateName)
		planOptions := []tfexec.PlanOption{
			tfexec.Out(planFile),
			tfexec.Refresh(true),
			tfexec.Parallelism(terraformParallelism),
		}

		// Then merge all tfvars found in config directory into -var-file options
		varOpts, err := terraform.ExpandVarDirectory(o.ConfigPath)
		cobra.CheckErr(err)
		planOptions = append(planOptions, varOpts...)

		// Now actually invoke Terraform plan
		console.StartSpinner()
		changes, err := tf.Plan(context.Background(), planOptions...)
		console.StopSpinner()
		cobra.CheckErr(err)
		if changes {
			console.Success("Plan contains infrastructure updates")
		} else {
			console.Success("Plan detected no changes")
			console.Success("Skipping the apply phase")
			console.Success("Rover completed")
			return
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
			tfexec.Parallelism(terraformParallelism),
		}

		// Now actually invoke Terraform apply
		console.StartSpinner()
		err := tf.Apply(context.Background(), applyOptions...)
		console.StopSpinner()
		cobra.CheckErr(err)

		// Special case for post launchpad deployment
		newStorageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
		cobra.CheckErr(err)
		if o.LaunchPadMode && existingStorageID != newStorageID {
			console.Info("Detected the launchpad infrastructure has been deployed or updated")

			stateFileName := o.OutPath + "/" + o.StateName + ".tfstate"
			azure.UploadFileToBlob(newStorageID, o.Workspace, o.StateName+".tfstate", stateFileName)
			console.Info("Uploading state from launchpad process to Azure storage")
			os.Remove(stateFileName)
		}

		console.Success("Apply was successful")
	}

	console.Success("Rover completed")
}

// Carry out Terraform init operation in launchpad mode has no backend state
func (o *Options) runLaunchpadInit(tf *tfexec.Terraform) error {
	console.Info("Running init for launchpad")

	console.StartSpinner()
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	console.StopSpinner()
	return err
}

// Carry out Terraform init operation with remote state backend
func (o *Options) runRemoteInit(tf *tfexec.Terraform, storageID string) error {
	console.Info("Running init with remote state")

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

	console.StartSpinner()
	err := tf.Init(context.Background(), initOptions...)
	cobra.CheckErr(err)
	console.StopSpinner()
	return err
}

// This function
func (o *Options) initializeCAF() *tfexec.Terraform {
	tfPath, err := terraform.Setup()
	cobra.CheckErr(err)

	os.Setenv("ARM_SUBSCRIPTION_ID", o.TargetSubscription)
	os.Setenv("ARM_TENANT_ID", o.Subscription.TenantID)
	os.Setenv("TF_VAR_tfstate_subscription_id", o.StateSubscription)
	os.Setenv("TF_VAR_tf_name", fmt.Sprintf("%s.tfstate", o.StateName))
	os.Setenv("TF_VAR_tf_plan", fmt.Sprintf("%s.tfplan", o.StateName))
	os.Setenv("TF_VAR_workspace", o.Workspace)
	os.Setenv("TF_VAR_level", o.Level)
	os.Setenv("TF_VAR_environment", o.CafEnvironment)
	os.Setenv("TF_VAR_rover_version", version.Value)
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
	localStatePath := fmt.Sprintf("%s/tfstates/%s/%s", os.Getenv("TF_DATA_DIR"), o.Level, o.Workspace)
	err = os.MkdirAll(localStatePath, os.ModePerm)
	cobra.CheckErr(err)
	o.OutPath = localStatePath

	// Create new TF exec with the working dir set to source
	tf, err := tfexec.NewTerraform(o.SourcePath, tfPath)
	cobra.CheckErr(err)

	// The debugging done here
	if console.DebugEnabled {
		// By default no output from Terraform is seen
		// Also note TF_LOG env var is ignored by tfexec
		tf.SetStdout(os.Stdout)
		tf.SetStderr(os.Stderr)

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
func (o *Options) cleanUp() {
	_ = os.Remove(o.SourcePath + "/backend.azurerm.tf")
	_ = os.Remove(o.OutPath + "/" + o.StateName + ".tfstate")
	_ = os.Remove(o.OutPath + "/" + o.StateName + ".tfplan")
	_ = os.Remove(os.Getenv("TF_DATA_DIR") + "/terraform.tfstate")
}

// By copying this file we enable teh azurerm backend and therefore remote state
func (o *Options) enableAzureBackend() {
	console.Info("Enabling backend state with backend.azurerm.tf file")
	err := utils.CopyFile(o.SourcePath+"/backend.azurerm", o.SourcePath+"/backend.azurerm.tf")
	cobra.CheckErr(err)
}
