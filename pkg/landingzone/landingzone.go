//
// Rover - Core execution of landingzone operations and actions
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"context"
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
const SecretTenantID = "tenant-id"
const SecretLowerSAName = "lower-storage-account-name"
const SecretLowerRGName = "lower-resource-group-name"

// Called by all CAF actions to set up Terraform and configure it for CAF landingzones
func (c *TerraformAction) prepareTerraformCAF(o *Options) *tfexec.Terraform {
	// Get current Azure details, subscription etc from CLI
	acct := azure.GetSubscription()
	ident := azure.GetIdentity()

	// If they weren't set already, fall back to logged in account subscription
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

	// Slight hack for now, we set debug on when in dry-run mode
	if o.DryRun {
		console.DebugEnabled = true
	}

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
		// By default no output from Terraform is seen, note TF_LOG env var is ignored by tfexec
		tf.SetStdout(os.Stdout)
		tf.SetStderr(os.Stderr)

		// This gives us some info level logs we can send to stdout
		tf.SetLogger(console.Printfer{})

		console.Debug("==== Execution Context ====")
		o.Debug()

		console.Debug("==== Environmental Variables ====")
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "ARM_") || strings.HasPrefix(env, "AZURE_") || strings.HasPrefix(env, "TF_") {
				console.Debugf("   %s\n", env)
			}
		}
	}

	// Remove old files, reset backend etc
	o.cleanUp()

	// Find state storage account for this environment and level
	c.launchPadStorageID, err = azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	if err != nil {
		if o.LaunchPadMode {
			console.Warning("No state storage account found, but running in launchpad mode, we can continue")
		} else {
			console.Errorf("No state storage account found for environment '%s' and level %d, please deploy a launchpad first!\n", o.CafEnvironment, o.Level)
			cobra.CheckErr("Can't deploy a landing zone without a launchpad")
		}
	} else {
		console.Infof("Located state storage account %s\n", c.launchPadStorageID)
	}
	return tf
}

// Runs init in the correct mode
func (c TerraformAction) runTerraformInit(o *Options, tf *tfexec.Terraform) {
	var err error
	if o.LaunchPadMode && c.launchPadStorageID == "" {
		err = o.runLaunchpadInit(tf, false)
	} else {
		err = o.runRemoteInit(tf, c.launchPadStorageID)
	}
	cobra.CheckErr(err)
}

// Carry out Terraform init operation in launchpad mode has no backend state
func (o *Options) runLaunchpadInit(tf *tfexec.Terraform, reconfigure bool) error {
	console.Info("Running init for launchpad")

	console.StartSpinner()
	// Validate that the indentity we are using is owner on subscription, not sure why but it's in rover v1 code
	isOwner, err := azure.CheckIsOwner(o.Identity.ObjectID, o.StateSubscription)
	cobra.CheckErr(err)
	if !isOwner {
		console.StopSpinner()
		console.Errorf("The identity %s (%s) is not assigned 'Owner' role on subscription %s\n", o.Identity.DisplayName, o.Identity.ObjectID, o.StateSubscription)
		cobra.CheckErr("To deploy a launchpad the identity used must be assigned the 'Owner' role")
	}

	// Proceed and run tf init
	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.Reconfigure(reconfigure))
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

	lpTenantID, err := kvClient.GetSecret(SecretTenantID)
	if err != nil {
		return err
	}
	lpLowerSAName, err := kvClient.GetSecret(SecretLowerSAName)
	if err != nil {
		return err
	}
	lpLowerResGrp, err := kvClient.GetSecret(SecretLowerRGName)
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
