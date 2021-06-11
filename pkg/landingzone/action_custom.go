package landingzone

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const customActionPath = "./custom_actions"

type CustomAction struct {
	ActionBase
	command CustomActionDefinition
}

type CustomActionDefinition struct {
	Name            string
	Executable      string
	Description     string
	ContinueOnError bool
	Arguments       []string
}

func NewCustomAction(cad CustomActionDefinition) CustomAction {
	return CustomAction{
		command: cad,
		ActionBase: ActionBase{
			name:        cad.Name,
			description: cad.Description,
		},
	}
}

func (c CustomAction) Execute(o *Options) error {
	console.Successf("Running custom action: %s %s\n", c.Name(), o.SourcePath)
	args := []string{}
	for _, argDefined := range c.command.Arguments {
		if argDefined == "{{SOURCE_DIR}}" {
			args = append(args, o.SourcePath)
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

func FetchCustomActions() ([]CustomAction, error) {
	actions := []CustomAction{}

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

		cad := CustomActionDefinition{}

		err = yaml.Unmarshal(buf, &cad)
		if err != nil {
			return nil, err
		}

		if cad.Name == "" {
			console.Warningf("Warning: custom action %s has no name it will be ignored\n", file.Name())
			continue
		}
		if cad.Executable == "" {
			console.Warningf("Warning: custom action %s has no executable it will be ignored\n", file.Name())
			continue
		}
		if cad.Description == "" {
			cad.Description = "No description"
		}

		actions = append(actions, NewCustomAction(cad))
	}

	return actions, nil
}
