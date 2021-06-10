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
	"path/filepath"
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
const secretTenantID = "tenant-id"
const secretLowerSAName = "lower-storage-account-name"
const secretLowerRGName = "lower-resource-group-name"

// Execute is entry point for `landingzone run`, `launchpad run` and `cd` operations
// This executes an action against a set of config options
func (o *Options) Execute(action Action) {
	console.Infof("Begin execution of action '%s'\n", action.Name())

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

	// All the env vars & other setup needed for CAF and get handle on Terraform
	tf := o.initializeCAF()

	// Remove old files, reset backend etc
	o.cleanUp()

	// Find state storage account for this environment and level
	existingStorageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	if err != nil {
		if o.LaunchPadMode {
			console.Warning("No state storage account found, but running in launchpad mode, we can continue")
		} else {
			console.Errorf("No state storage account found for environment '%s' and level %d, please deploy a launchpad first!\n", o.CafEnvironment, o.Level)
			cobra.CheckErr("Can't deploy a landing zone without a launchpad")
		}
	} else {
		console.Infof("Located state storage account %s\n", existingStorageID)
	}

	// TODO: PUT COMMMANDS HERE THAT DONT NEED INIT AND EXIT EARLY
	if action == ActionFormat {
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
		err = o.runLaunchpadInit(tf, false)
	} else {
		err = o.runRemoteInit(tf, existingStorageID)
	}
	cobra.CheckErr(err)

	// If the action is just init, then stop here and don't proceed
	if action == ActionInit {
		console.Success("Rover completed, only init was run and no infrastructure changes were planned or applied")
		return
	}

	console.Success("Init completed, moving to next phase")

	//
	// Terraform plan step
	//
	planChanges := false
	if action == ActionPlan || action == ActionApply {
		console.Info("Carrying out the Terraform plan phase")

		// Connect to launchpad, setting all the vars needed by the landingzone
		if !o.LaunchPadMode {
			err = o.connectToLaunchPad(existingStorageID)
			cobra.CheckErr(err)
		}

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
		for _, vo := range varOpts {
			// Note. spread operator would not work here, I tried ¯\_(ツ)_/¯
			planOptions = append(planOptions, vo)
		}

		console.StartSpinner()
		planChanges, err = tf.Plan(context.Background(), planOptions...)
		console.StopSpinner()
		cobra.CheckErr(err)
		if planChanges {
			console.Success("Plan contains infrastructure updates")
		} else {
			console.Success("Plan detected no changes")
			console.Success("Any apply step will be skipped")
		}
	}

	//
	// Terraform apply step, won't run if plan found no changes
	//
	if action == ActionApply && planChanges {
		console.Info("Carrying out the Terraform apply phase")

		planFile := fmt.Sprintf("%s/%s.tfplan", o.OutPath, o.StateName)
		stateFile := fmt.Sprintf("%s/%s.tfstate", o.OutPath, o.StateName)

		// Build apply options, with plan file and state out
		applyOptions := []tfexec.ApplyOption{
			tfexec.DirOrPlan(planFile),
			tfexec.StateOut(stateFile),
			tfexec.Parallelism(terraformParallelism),
		}

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

	if action == ActionValidate {
		console.Info("Carrying out the Terraform validate phase")

		// Now actually invoke Terraform apply
		console.StartSpinner()
		_, err := tf.Validate(context.Background())
		console.StopSpinner()
		cobra.CheckErr(err)
	}

	//
	// Terraform validate step
	//
	if action == ActionValidate {
		console.Info("Carrying out the Terraform validate phase")

		console.StartSpinner()
		_, err := tf.Validate(context.Background())
		console.StopSpinner()
		cobra.CheckErr(err)
	}

	//
	// Destroy action
	//
	if action == ActionDestroy {
		console.Info("Carrying out the Terraform destroy phase")

		stateFileName := o.OutPath + "/" + o.StateName + ".tfstate"

		// Build apply options, with plan file and state out
		destroyOptions := []tfexec.DestroyOption{
			tfexec.Parallelism(terraformParallelism),
			tfexec.Refresh(false),
		}

		// We need to do all sorts of extra shenanigans for launchPadMode
		if o.LaunchPadMode {
			console.Warning("WARNING! You are destroying the launchpad!")
			if existingStorageID == "" {
				console.Error("Looks like this launchpad has already been deleted, bye!")
				cobra.CheckErr("Destroy was aborted")
			}

			// IMPORTANT!
			o.cleanUp()

			// Download the current state
			azure.DownloadFileFromBlob(existingStorageID, o.Workspace, o.StateName+".tfstate", stateFileName)

			// Reset back to use local state
			console.Warning("Resetting state to local, have to re-run init")
			err = o.runLaunchpadInit(tf, true)
			cobra.CheckErr(err)
			// IMPORTANT!
			_ = os.Remove(o.SourcePath + "/backend.azurerm.tf")

			// Tell destroy to use local downloaded state to destroy a launchpad
			// TODO: This is a deprecated option, the solution is to switch to this
			// https://www.terraform.io/docs/language/settings/backends/local.html
			// nolint
			destroyOptions = append(destroyOptions, tfexec.State(stateFileName))
		} else {
			// Connect to launchpad, setting all the vars needed by the landingzone
			err = o.connectToLaunchPad(existingStorageID)
			cobra.CheckErr(err)
		}

		// Merge all tfvars found in config directory into -var-file options
		varOpts, err := terraform.ExpandVarDirectory(o.ConfigPath)
		cobra.CheckErr(err)
		for _, vo := range varOpts {
			// Note. spread operator would not work here, I tried ¯\_(ツ)_/¯
			destroyOptions = append(destroyOptions, vo)
		}

		console.Warning("Destroy is now running ...")
		console.StartSpinner()
		err = tf.Destroy(context.Background(), destroyOptions...)
		console.StopSpinner()
		cobra.CheckErr(err)

		// Remove files
		o.cleanUp()
		_ = os.RemoveAll(o.OutPath)
		console.Success("Destroy was successful")
	}
	console.Successf("Execution of '%s' completed\n", action.Name())
}

// Carry out Terraform init operation in launchpad mode has no backend state
func (o *Options) runLaunchpadInit(tf *tfexec.Terraform, reconfigure bool) error {
	console.Info("Running init for launchpad")

	console.StartSpinner()
	err := tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.Reconfigure(reconfigure))
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
	filepath.Join(o.SourcePath, "backend.azurerm")
	err := utils.CopyFile(filepath.Join(o.SourcePath, "backend.azurerm"), filepath.Join(o.SourcePath, "backend.azurerm.tf"))
	cobra.CheckErr(err)
}

// Sets various TF_VAR_ variables required for a landingzone to be deployed/destroyed
func (o *Options) connectToLaunchPad(lpStorageID string) error {
	console.Infof("Connecting to launchpad for level '%s'\n", o.Level)
	lpKeyVaultID, err := azure.FindKeyVault(o.Level, o.CafEnvironment, o.StateSubscription)
	if err != nil {
		return err
	}
	if lpKeyVaultID == "" {
		return fmt.Errorf("Unable to locate the launchpad for environment '%s' and level '%s'", o.CafEnvironment, o.Level)
	}

	_, _, keyVaultName := azure.ParseResourceID(lpKeyVaultID)

	kvClient, err := azure.NewKVClient(azure.KeyvaultEndpointForSubscription(), keyVaultName)
	if err != nil {
		return err
	}

	lpTenantID, err := kvClient.GetSecret(secretTenantID)
	if err != nil {
		return err
	}
	lpLowerSAName, err := kvClient.GetSecret(secretLowerSAName)
	if err != nil {
		return err
	}
	lpLowerResGrp, err := kvClient.GetSecret(secretLowerRGName)
	if err != nil {
		return err
	}

	if lpLowerSAName == "" || lpTenantID == "" || lpLowerResGrp == "" {
		return fmt.Errorf("Required secret(s) not found in launchpad, either you are not authorized or the launchpad was not deployed correctly")
	}

	_, lpStorageResGrp, lpStorageName := azure.ParseResourceID(lpStorageID)

	console.Success("Connected to launchpad OK")
	console.Debugf(" - TF_VAR_tenant_id=%s\n", lpTenantID)
	console.Debugf(" - TF_VAR_tfstate_storage_account_name=%s\n", lpStorageName)
	console.Debugf(" - TF_VAR_tfstate_resource_group_name=%s\n", lpStorageResGrp)
	console.Debugf(" - TF_VAR_lower_storage_account_name=%s\n", lpLowerSAName)
	console.Debugf(" - TF_VAR_lower_resource_group_name=%s\n", lpLowerResGrp)

	_ = os.Setenv("TF_VAR_tenant_id", lpTenantID)
	_ = os.Setenv("TF_VAR_tfstate_storage_account_name", lpStorageName)
	_ = os.Setenv("TF_VAR_tfstate_resource_group_name", lpStorageResGrp)
	_ = os.Setenv("TF_VAR_lower_storage_account_name", lpLowerSAName)
	_ = os.Setenv("TF_VAR_lower_resource_group_name", lpLowerResGrp)

	_ = os.Setenv("TF_VAR_tfstate_container_name", o.Workspace)
	_ = os.Setenv("TF_VAR_lower_container_name", o.Workspace)
	// NOTE: This will have been set by initializeCAF()
	_ = os.Setenv("TF_VAR_tfstate_key", os.Getenv("TF_VAR_tf_name"))

	return nil
}
