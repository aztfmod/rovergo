package custom

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const customActionPath = "./custom_actions"

// Action is an custom action implementation which runs external executables
type Action struct {
	landingzone.ActionBase
	command actionDefinition
}

// customActionDefinition is used to parse the YAML config files
type actionDefinition struct {
	Name            string
	Executable      string
	Description     string
	ContinueOnError bool
	Arguments       []string
}

// This is never called externally, only by calling FetchCustomActions
func newCustomAction(cad actionDefinition) Action {
	return Action{
		command: cad,
		ActionBase: landingzone.ActionBase{
			Name:        cad.Name,
			Description: cad.Description,
		},
	}
}

// Execute runs this custom action by running the external executable
func (c Action) Execute(o *landingzone.Options) error {
	console.Successf("Running custom action: %s %s\n", c.Name, o.SourcePath)
	args := []string{}

	// This implements a simple variable subsistution syntax
	for _, argDefined := range c.command.Arguments {
		if argDefined == "{{SOURCE_DIR}}" {
			args = append(args, o.SourcePath)
			continue
		}
		if argDefined == "{{CONFIG_DIR}}" {
			args = append(args, o.ConfigPath)
			continue
		}
		if argDefined == "{{LEVEL}}" {
			args = append(args, o.Level)
			continue
		}
		if argDefined == "{{STATE_NAME}}" {
			args = append(args, o.StateName)
			continue
		}
		if argDefined == "{{CAF_ENV}}" {
			args = append(args, o.CafEnvironment)
			continue
		}
		if argDefined == "{{WORKSPACE}}" {
			args = append(args, o.Workspace)
			continue
		}
		args = append(args, argDefined)
	}

	cmd := command.NewCommand(c.command.Executable, args)
	cmd.Silent = false
	err := cmd.Execute()

	console.Error(cmd.StdErr)
	console.Success(cmd.StdOut)

	// NOTE: When running across multiple levels/stacks
	// We will exit early when we hit first error, this could be improved
	cobra.CheckErr(err)

	return nil
}

// FetchActions is called by root cmd during init
// It finds all the custom action defintions and returns them to be plugged into the CLI
func FetchActions() ([]landingzone.Action, error) {
	actions := []landingzone.Action{}

	// Finds all .tfvars in directory, note. we no longer use walk as it was recursive
	actionFiles, err := os.ReadDir(customActionPath)
	if err != nil {
		console.Warning("Warning: No rover custom_actions directory found")
		return nil, nil
	}

	for _, file := range actionFiles {
		if !(strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
			continue
		}

		buf, err := os.ReadFile(filepath.Join(customActionPath, file.Name()))
		if err != nil {
			return nil, err
		}

		definition := actionDefinition{}

		err = yaml.Unmarshal(buf, &definition)
		if err != nil {
			return nil, err
		}

		if definition.Name == "" {
			console.Warningf("Warning: custom action %s has no name it will be ignored\n", file.Name())
			continue
		}
		if definition.Executable == "" {
			console.Warningf("Warning: custom action %s has no executable it will be ignored\n", file.Name())
			continue
		}
		if definition.Description == "" {
			definition.Description = "No description"
		}

		actions = append(actions, newCustomAction(definition))
	}

	return actions, nil
}
