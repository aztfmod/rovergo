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
)

var requiredMinVer, _ = version.NewVersion("0.15.0")

// Setup gets a path to Terraform, optionally install it, and check version
func Setup() (string, error) {

	// First look in system path & installPath
	path, err := tfinstall.Find(context.Background(), tfinstall.LookPath())

	if err != nil {
		console.Error("Failed to find Terraform installation")
		os.Exit(1)
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

// ExpandVarDirectory returns an array of var file options from a directory of tfvars
func ExpandVarDirectory(varDir string) ([]*tfexec.VarFileOption, error) {
	varFileOpts := []*tfexec.VarFileOption{}

	// Finds all .tfvars in directory, note. we no longer use walk as it was recursive
	tfvarFiles, err := os.ReadDir(varDir)
	if err != nil {
		return nil, err
	}
	for _, file := range tfvarFiles {
		if !strings.HasSuffix(file.Name(), ".tfvars") {
			continue
		}
		varFileName := filepath.Join(varDir, file.Name())
		varFileOpts = append(varFileOpts, tfexec.VarFile(varFileName))
		console.Debugf("Found tfvar file to use: %s\n", varFileName)
	}

	// Ensure we have some tfvars, otherwise we're going to have a really bad time
	if len(varFileOpts) <= 0 {
		return nil, fmt.Errorf("failed to find any tfvars files in config directory: %s", varDir)
	}

	return varFileOpts, nil
}
