//
// Wrapper and helper for running external commands
//

package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
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
}

func NewCommand(exe string, args []string) *Command {
	return &Command{
		Exe:    exe,
		Args:   args,
		DryRun: false,
	}
}

func (c *Command) Execute(includeOsEnv bool) error {
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
	color.Blue("Executing %s %s", c.Exe, c.Args)

	// Handy for debugging
	if c.DryRun {
		return nil
	}

	// Append system env vars, pretty rare you *wouldn't* want these
	if includeOsEnv {
		cmd.Env = append(cmd.Env, os.Environ()...)
	}

	// Set buffers to capture stdout & stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Actually run the thing
	if err := cmd.Run(); err != nil {
		color.Red("Failed, %s", stderr.String())
		return err
	}

	c.StdOut = stdout.String()
	c.StdErr = stderr.String()
	return nil
}

func EnsureDirectory(dir string) {
	os.MkdirAll(dir, os.ModePerm)
}

func RemoveDirectory(dir string) {
	os.RemoveAll(dir)
}

func CheckCommand(reqCmdName string) error {
	_, err := exec.LookPath(reqCmdName)
	if err != nil {
		return fmt.Errorf(color.RedString("required command %s not found in system path", reqCmdName))
	}
	return nil
}
