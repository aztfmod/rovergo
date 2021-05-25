//
// Wrapper for running terraform commands and interating with Terraform
//

package terraform

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var requiredMinVer, _ = version.NewVersion("0.15.0")

//
// Setup gets a path to Terraform, optionally install it, and check version
//
func Setup() (string, error) {
	// Config to control if install happens and where
	install := viper.GetBool("terraform.install")
	installPath := viper.GetString("terraform.install-path")

	// First look in system path & installPath
	path, err := tfinstall.Find(context.Background(), tfinstall.LookPath(), tfinstall.ExactPath(installPath+"/terraform"))

	// Try to install and then locate terraform
	if err != nil && install {
		color.Yellow("Attempting install of terraform into %s", installPath)
		// Any error from install is lost and never set, probably a bug
		_, _ = tfinstall.Find(context.Background(), tfinstall.LatestVersion(installPath, false))
		path, err = tfinstall.Find(context.Background(), tfinstall.ExactPath(installPath+"/terraform"))
		if err != nil {
			color.Red("Install failed, make sure %s exists and is writable", installPath)
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	// Initialize terraform for use with Azure
	SetAzureEnvVars()
	CheckVersion(path)

	path, err = filepath.Abs(path)
	cobra.CheckErr(err)

	return path, nil
}

//
// CheckVersion ensures that Terraform is at the required version
//
func CheckVersion(path string) {
	// Working dir shouldn't matter for this command
	tf, err := tfexec.NewTerraform(".", path)
	cobra.CheckErr(err)
	tfVer, _, err := tf.Version(context.Background(), false)
	cobra.CheckErr(err)

	if !tfVer.GreaterThanOrEqual(requiredMinVer) {
		cobra.CheckErr(color.RedString("Terrform version %v is behind required minimum %v", tfVer, requiredMinVer))
	}
	color.Green("Terraform is at version %v", tfVer)
}

//
// SetupAzureEnvironment should be called before any terraform operations
//
func SetAzureEnvVars() {
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
