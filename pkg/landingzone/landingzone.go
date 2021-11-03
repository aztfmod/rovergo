//
// Rover - Core execution of landingzone operations and actions
// * Curent status is: launchpad deploy works and sets up remote state
// * Ben C, May 2021
//

package landingzone

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
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
func (c *TerraformAction) prepareTerraformCAF(o *Options) (*tfexec.Terraform, error) {

	err := o.SetupEnvironment()
	if err != nil {
		return nil, err
	}

	// Locate terraform
	tfPath, err := terraform.Setup()
	if err != nil {
		return nil, err
	}

	// Create new TF exec with the working dir set to source
	tf, err := tfexec.NewTerraform(o.SourcePath, tfPath)
	if err != nil {
		return nil, err
	}

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
			console.Errorf("No state storage account found for environment '%s' and level %s, please deploy a launchpad first!\n", o.CafEnvironment, o.Level)
			return nil, errors.New("can't deploy a landing zone without a launchpad")
		}
	} else {
		console.Infof("Located state storage account %s\n", c.launchPadStorageID)
	}

	return tf, nil
}

// SetupEnvironment for all the terraform env vars AND values in options stuct
func (o *Options) SetupEnvironment() error {
	// Get current Azure details, subscription etc from CLI
	acct, err := azure.GetSubscription()
	cobra.CheckErr(err)

	// If they weren't set already, fall back to logged in account subscription
	if o.StateSubscription == "" {
		o.StateSubscription = acct.ID
	}
	if o.TargetSubscription == "" {
		o.TargetSubscription = acct.ID
	}
	o.Subscription = *acct

	if o.LaunchPadMode {
		if o.TargetSubscription != o.StateSubscription {
			return errors.New("in launchpad mode, state-sub and target-sub must be the same Azure subscription")
		}
	}

	// Get the currently signed in identity regardless of type
	o.Identity = getIdentity(*acct, o.TargetSubscription)
	console.Successf("Obtained identity successfully.\nWe are signed in as: %s '%s' (%s)\n", o.Identity.ObjectType, o.Identity.DisplayName, o.Identity.ObjectID)

	// Slight hack for now, we set debug on when in dry-run mode
	if o.DryRun {
		console.DebugEnabled = true
	}

	if strings.EqualFold(o.Identity.ObjectType, "servicePrincipal") {
		os.Setenv("ARM_CLIENT_ID", o.Identity.ClientID)

		// the AssignedIdentityInfo starts with "MSI" for both system assigned and user assigned
		if strings.HasPrefix(acct.User.AssignedIdentityInfo, "MSI") {
			os.Setenv("ARM_USE_MSI", "true")
		} else {
			// Otherwise were using a old fashioned SP and we need the secret to be set outside of rover
			if os.Getenv("ARM_CLIENT_SECRET") == "" && os.Getenv("ARM_CLIENT_CERTIFICATE_PATH") == "" {
				return errors.New("when signed in as service principal, you must set ARM_CLIENT_SECRET or ARM_CLIENT_CERTIFICATE_PATH")
			}
		}
	}

	os.Setenv("TF_DATA_DIR", o.DataDir)
	os.Setenv("TF_VAR_tfstate_key", "")
	os.Setenv("ARM_SUBSCRIPTION_ID", o.TargetSubscription)
	os.Setenv("ARM_TENANT_ID", o.Subscription.TenantID)
	os.Setenv("TF_VAR_tfstate_subscription_id", o.StateSubscription)
	os.Setenv("TF_VAR_tf_name", path.Join(o.StateName, ".tfstate"))
	os.Setenv("TF_VAR_tf_plan", path.Join(o.StateName, ".tfplan"))
	os.Setenv("TF_VAR_workspace", o.Workspace)
	os.Setenv("TF_VAR_level", o.Level)
	os.Setenv("TF_VAR_environment", o.CafEnvironment)
	os.Setenv("TF_VAR_rover_version", version.Value)
	os.Setenv("TF_VAR_tenant_id", o.Subscription.TenantID)
	os.Setenv("TF_VAR_user_type", o.Identity.ObjectType)
	os.Setenv("TF_VAR_logged_user_objectId", o.Identity.ObjectID)

	return nil
}

// Try to get our identity which might be user, managed-identity or service principal
func getIdentity(acct azure.Subscription, targetSubID string) azure.Identity {
	if strings.EqualFold(acct.User.Usertype, "user") {
		console.Debug("Detected we are signed in as a user. Attempting to get identity from CLI")
		ident, err := azure.GetSignedInIdentity()
		cobra.CheckErr(err)
		return *ident

	} else if strings.HasPrefix(acct.User.AssignedIdentityInfo, "MSI") {
		console.Debug("Detected we are signed in as MSI. Attempting to get VM assigned identity")

		userAssignedByObjectID := strings.HasPrefix(acct.User.AssignedIdentityInfo, "MSIObject")
		userAssignedByClientID := strings.HasPrefix(acct.User.AssignedIdentityInfo, "MSIClient")
		systemAssigned := (acct.User.AssignedIdentityInfo == "MSI")

		var vmIdentityID string
		if userAssignedByObjectID || userAssignedByClientID {
			vmIdentityID = strings.SplitAfterN(acct.User.AssignedIdentityInfo, "-", 2)[1]
		}

		metadata := azure.VMInstanceMetadataService()
		vmIdentities, err := azure.GetVMIdentities(acct.ID, metadata.Compute.ResourceGroupName, metadata.Compute.Name)
		cobra.CheckErr(err)

		// look for the vm identity that matches the az login id
		// it could be a system assigned (AssignedIdentityInfo="MSI")
		// it could be a user assigned (AssignedIdentityInfo="MSIObject" or "MSIClient")
		for _, id := range vmIdentities {

			if systemAssigned {
				if id.DisplayName == "SystemAssigned" {
					return id
				}
			} else if (userAssignedByObjectID && id.ObjectID == vmIdentityID) || (userAssignedByClientID && id.ClientID == vmIdentityID) {
				return id
			}
		}

		return azure.Identity{}

	} else if strings.EqualFold(acct.User.Usertype, "serviceprincipal") {
		console.Debug("Detected we are signed in as a service principal. Attempting to get identity from the Graph API")
		// The Azure CLI puts the SP clientid in the name field, which is weird but useful for us
		identity, err := azure.GetServicePrincipalIdentity(acct.User.Name)
		cobra.CheckErr(err)
		return *identity
	} else {
		console.Error("Signed in identity is of unknown type")
		console.Errorf("%+v", acct)
		cobra.CheckErr("Rover cannot continue")
	}

	// We never get here, but go compiler doesn't understand that
	return azure.Identity{}
}

// Runs init in the correct mode
func (c TerraformAction) runTerraformInit(o *Options, tf *tfexec.Terraform, forceLocal bool) {
	var err error
	o.removeStateConfig()

	if (o.LaunchPadMode && c.launchPadStorageID == "") || forceLocal {
		err = o.runLaunchpadInit(tf, false)
	} else {
		err = o.runRemoteInit(tf, c.launchPadStorageID)
	}
	cobra.CheckErr(err)
}

// Carry out Terraform init operation in launchpad mode has no backend state
func (o *Options) runLaunchpadInit(tf *tfexec.Terraform, reconfigure bool) error {
	console.Info("Running init for launchpad (local state)")

	console.StartSpinner()
	// Validate that the identity we are using is owner on subscription, not sure why but it's in rover v1 code
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

	subID, resGrp, accountName, err := azure.ParseResourceID(storageID)
	if err != nil {
		return err
	}
	accessKey, err := azure.GetAccountKey(subID, accountName, resGrp)
	cobra.CheckErr(err)

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
	err = tf.Init(context.Background(), initOptions...)
	cobra.CheckErr(err)
	console.StopSpinner()
	return err
}

// Remove files to ensure a clean run
// TODO: This may require future changes please leave the commented out lines
func (o *Options) cleanUp() {
	//_ = os.Remove(o.SourcePath + "/backend.azurerm.tf")
	_ = os.Remove(o.DataDir + "/" + o.StateName + ".tfstate")
	//_ = os.Remove(o.DataDir + "/terraform.tfstate")
}

// Remove the remote state configuration
// TODO: This may require future changes please leave the commented out lines
func (o *Options) removeStateConfig() {
	_ = os.Remove(o.SourcePath + "/backend.azurerm.tf")
	_ = os.Remove(o.DataDir + "/terraform.tfstate")
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

	_, _, keyVaultName, err := azure.ParseResourceID(lpKeyVaultID)
	if err != nil {
		return err
	}

	endpoint, err := azure.KeyvaultEndpointForSubscription()
	if err != nil {
		return err
	}
	kvClient, err := azure.NewKVClient(endpoint, keyVaultName)
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

	_, lpStorageResGrp, lpStorageName, err := azure.ParseResourceID(lpStorageID)
	if err != nil {
		return err
	}

	console.Success("Connected to launchpad OK")
	_ = os.Setenv("TF_VAR_tenant_id", lpTenantID)
	_ = os.Setenv("TF_VAR_tfstate_storage_account_name", lpStorageName)
	_ = os.Setenv("TF_VAR_tfstate_resource_group_name", lpStorageResGrp)
	_ = os.Setenv("TF_VAR_lower_storage_account_name", lpLowerSAName)
	_ = os.Setenv("TF_VAR_lower_resource_group_name", lpLowerResGrp)

	_ = os.Setenv("TF_VAR_tfstate_container_name", o.Workspace)
	_ = os.Setenv("TF_VAR_lower_container_name", o.Workspace)
	// NOTE: This will have been set by initializeCAF()
	//_ = os.Setenv("TF_VAR_tfstate_key", os.Getenv("TF_VAR_tf_name"))
	_ = os.Setenv("TF_VAR_tfstate_key", fmt.Sprintf("%s.tfstate", o.StateName))

	console.Debugf(" - TF_VAR_tenant_id=%s\n", lpTenantID)
	console.Debugf(" - TF_VAR_tfstate_storage_account_name=%s\n", lpStorageName)
	console.Debugf(" - TF_VAR_tfstate_resource_group_name=%s\n", lpStorageResGrp)
	console.Debugf(" - TF_VAR_lower_storage_account_name=%s\n", lpLowerSAName)
	console.Debugf(" - TF_VAR_lower_resource_group_name=%s\n", lpLowerResGrp)
	console.Debugf(" - TF_VAR_tfstate_key=%s\n", os.Getenv("TF_VAR_tfstate_key"))

	return nil
}
