package custom

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//go:embed default_actions.yaml
var defaultFileContent string

const actionsFilename = "actions.yaml"

// Action is an custom action implementation which runs external executables
type Action struct {
	landingzone.ActionBase
	command actionDefinition
}

// actionDefinition is used to parse the YAML config files
type actionDefinition struct {
	Executable  string
	Description string
	Arguments   []string
}

// This is never called externally, only by calling FetchCustomActions
func newCustomAction(name string, cad actionDefinition) Action {
	return Action{
		command: cad,
		ActionBase: landingzone.ActionBase{
			Name:        name,
			Description: cad.Description + " [Custom]",
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
func FetchActions() (actions []landingzone.Action, err error) {
	roverHomeDir, _ := utils.GetRoverDirectory()
	custActionsPath := filepath.Join(roverHomeDir, actionsFilename)
	_, err = os.Stat(custActionsPath)
	// If doesn't exist then place our default YAML file in .rover
	if err != nil {
		fileErr := ioutil.WriteFile(custActionsPath, []byte(defaultFileContent), 0777)
		if fileErr != nil {
			return nil, fileErr
		}
	}

	// Read file and unmarshall
	file, err := os.Open(custActionsPath)
	if err != nil {
		return nil, err
	}
	// The actions YAML file is a map of strings to definitions, where the key is the name of the action
	actionsYaml := map[string]actionDefinition{}
	decoder := yaml.NewDecoder(file)
	// Enabling strict mode prevents duplicate keys
	decoder.SetStrict(true)
	err = decoder.Decode(&actionsYaml)
	if err != nil {
		return nil, err
	}

	for actionName, actionDef := range actionsYaml {
		actions = append(actions, newCustomAction(actionName, actionDef))
	}

	return
}
