//
// Rover - Terraform helper
// * To assist calling tf-exec for running Terrafrom CLI
// * Ben C, May 2021
//

package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var requiredMinVer, _ = version.NewVersion("0.15.0")

// Setup gets a path to Terraform, optionally install it, and check version
func Setup() (string, error) {
	// Config to control if install happens and where
	install := viper.GetBool("terraform.install")
	installPath := viper.GetString("terraform.install-path")

	// First look in system path & installPath
	path, err := tfinstall.Find(context.Background(), tfinstall.LookPath(), tfinstall.ExactPath(installPath+"/terraform"))

	// Try to install and then locate terraform
	if err != nil && install {
		console.Infof("Attempting install of terraform into %s\n", installPath)
		// Any error from install is lost and never set, probably a bug
		_, _ = tfinstall.Find(context.Background(), tfinstall.LatestVersion(installPath, false))
		path, err = tfinstall.Find(context.Background(), tfinstall.ExactPath(installPath+"/terraform"))
		if err != nil {
			console.Errorf("Install failed, make sure %s exists and is writable\n", installPath)
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	// Initialize terraform for use with Azure
	SetEnvVars()
	CheckVersion(path)

	path, err = filepath.Abs(path)
	cobra.CheckErr(err)

	return path, nil
}

// CheckVersion ensures that Terraform is at the required version
func CheckVersion(path string) {
	// Working dir shouldn't matter for this command
	tf, err := tfexec.NewTerraform(".", path)
	cobra.CheckErr(err)
	tfVer, _, err := tf.Version(context.Background(), false)
	cobra.CheckErr(err)

	if !tfVer.GreaterThanOrEqual(requiredMinVer) {
		cobra.CheckErr(fmt.Sprintf("Terrform version %v is behind required minimum %v", tfVer, requiredMinVer))
	}
	console.Successf("Terraform is at version %v\n", tfVer)
}

// SetEnvVars should be called before any terraform operations
// It essentially "logs in" to Terraform with the creds stored in config
// TODO: REMOVE DEPRECATED 🔥
func SetEnvVars() {
	os.Setenv("ARM_SUBSCRIPTION_ID", viper.GetString("auth.subscription-id"))
	os.Setenv("ARM_CLIENT_ID", viper.GetString("auth.client-id"))
	os.Setenv("ARM_TENANT_ID", viper.GetString("auth.tenant-id"))
	os.Setenv("ARM_ENVIRONMENT", viper.GetString("auth.environment"))
	os.Setenv("ARM_CLIENT_CERTIFICATE_PATH", viper.GetString("auth.client-cert-path"))
	os.Setenv("ARM_CLIENT_CERTIFICATE_PASSWORD", viper.GetString("auth.client-cert-password"))
	os.Setenv("ARM_CLIENT_SECRET", viper.GetString("auth.client-secret"))
	os.Setenv("ARM_USE_MSI", viper.GetString("auth.use-msi"))
	os.Setenv("ARM_MSI_ENDPOINT", viper.GetString("auth.msi-endpoint"))
}

// Authenticate will attempt to auth using the go-azure-helper and return auth config
// TODO: REMOVE DEPRECATED 🔥
func Authenticate() (*authentication.Config, error) {
	builder := &authentication.Builder{
		TenantOnly:                     false,
		SupportsAuxiliaryTenants:       false,
		AuxiliaryTenantIDs:             nil,
		SupportsAzureCliToken:          true,
		ClientID:                       viper.GetString("auth.client-id"),
		ClientSecret:                   viper.GetString("auth.client-secret"),
		SubscriptionID:                 viper.GetString("auth.subscription-id"),
		TenantID:                       viper.GetString("auth.tenant-id"),
		SupportsClientCertAuth:         true,
		SupportsClientSecretAuth:       true,
		SupportsManagedServiceIdentity: viper.GetBool("auth.use-msi"), // TODO: Should this be conditional on ARM_USE_MSI ?
	}

	return builder.Build()
}
