//
// Rover - Command line wrapper
// * Helper functions for calling os/exec in a standard way
// * Ben C, May 2021
//

package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
)

type EnvVar struct {
	Name  string
	Value string
}

type Command struct {
	Exe     string
	Args    []string
	EnvVars []EnvVar
	StdErr  string
	StdOut  string
	DryRun  bool
	Silent  bool
	OsEnv   bool
}

func NewCommand(exe string, args []string) *Command {
	return &Command{
		Exe:    exe,
		Args:   args,
		DryRun: false,
		Silent: true,
		OsEnv:  true,
	}
}

func (c *Command) Execute() error {
	if err := CheckCommand(c.Exe); err != nil {
		return err
	}

	cmd := exec.Command(c.Exe, c.Args...)
	// Set extra env vars if they exist
	for _, envVar := range c.EnvVars {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("%s=%s", envVar.Name, envVar.Value),
		)
	}

	if !c.Silent {
		console.Debugf("Executing %s %s\n", c.Exe, c.Args)
	}

	// Handy for debugging
	if c.DryRun {
		return nil
	}

	// Append system env vars, pretty rare you *wouldn't* want these
	if c.OsEnv {
		cmd.Env = append(cmd.Env, os.Environ()...)
	}

	// Set buffers to capture stdout & stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Actually run the thing
	if err := cmd.Run(); err != nil {
		console.Errorf("Failed, %s", stderr.String())
		c.StdOut = stdout.String()
		c.StdErr = stderr.String()
		return err
	}

	c.StdOut = stdout.String()
	c.StdErr = stderr.String()
	return nil
}

func QuickRun(args ...string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("must supply at least one argument")
	}
	exe := args[0]
	restOfArgs := utils.StringSliceDel(args, 0)
	cmd := NewCommand(exe, restOfArgs)
	err := cmd.Execute()
	if err != nil {
		return "", err
	}
	return cmd.StdOut, nil
}

func EnsureDirectory(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	cobra.CheckErr(err)
}

func RemoveDirectory(dir string) {
	err := os.RemoveAll(dir)
	cobra.CheckErr(err)
}

func CheckCommand(reqCmdName string) error {
	_, err := exec.LookPath(reqCmdName)
	if err != nil {
		return fmt.Errorf("required command %s not found in system path", reqCmdName)
	}
	return nil
}

func ValidateDependencies() {
	azErr := CheckCommand("az")
	if azErr != nil {
		console.Errorf("The %s.\nPlease install from https://docs.microsoft.com/en-us/cli/azure/install-azure-cli", azErr.Error())
	}

	tfErr := CheckCommand("terraform")
	if tfErr != nil {
		console.Errorf("The %s.\nPlease install from https://www.terraform.io/downloads.html \n", tfErr.Error())
	}

	if azErr != nil || tfErr != nil {
		os.Exit(1)
	}
}
