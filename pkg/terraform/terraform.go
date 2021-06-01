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
	"strings"

	"github.com/aztfmod/rover/pkg/console"
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

// ExpandVarDirectory returns an array of plan options from a directory of tfvars
func ExpandVarDirectory(varDir string) ([]tfexec.PlanOption, error) {
	planOptions := []tfexec.PlanOption{}

	// Finds all .tfvars in directory and recurse down
	err := filepath.Walk(varDir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".tfvars") {
			return nil
		}

		planOptions = append(planOptions, tfexec.VarFile(path))
		console.Debugf("Found var file to use: %s\n", path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return planOptions, nil
}
