package custom

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aztfmod/rover/pkg/builtin/actions"
	commandpkg "github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const commandsFileName = "commands.yml"

// Action is an custom command implementation which runs external executables
type Action struct {
	landingzone.ActionBase
	Commands []Command
}

type yamlDefinition struct {
	Commands map[string]Command  `yaml:"commands"`
	Groups   map[string][]string `yaml:"groups"`
}

type Command struct {
	ExecutableName string `yaml:"executableName"`
	SubCommand     string `yaml:"subCommand"`
	Flags          string `yaml:"flags"`
	Debug          bool   `yaml:"debug"`
	RequiresInit   bool   `yaml:"requiresInit"`
	SetupEnv       bool   `yaml:"setupEnv"`
	Parameters     []struct {
		Name   string `yaml:"name"`
		Value  string `yaml:"value"`
		Prefix string `yaml:"prefix"`
	} `yaml:"parameters"`
}

func InitializeCustomCommandsAndGroups() error {
	commands, err := LoadCustomCommandsAndGroups()
	if err != nil {
		console.Errorf("Loading custom commands and groups failed: %s\n", err)
		return err
	}
	for _, ca := range commands {
		actions.ActionMap[ca.GetName()] = ca
	}
	return nil
}

// LoadCustomCommandsAndGroups is called by root cmd during init
// It finds all the custom action definitions and returns them to be plugged into the CLI
func LoadCustomCommandsAndGroups() (commands []landingzone.Action, err error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	commandsFilePath := filepath.Join(currentWorkingDirectory, commandsFileName)

	commandsFileContent, err := utils.ReadYamlFile(commandsFilePath)
	if err != nil {
		roverHomeDir, err := rover.HomeDirectory()
		if err != nil {
			return nil, err
		}
		commandsFilePath = filepath.Join(roverHomeDir, commandsFileName)

		commandsFileContent, err = utils.ReadYamlFile(commandsFilePath)
		if err != nil {
			return nil, err
		}
	}

	if len(commandsFileContent) == 0 {
		return nil, fmt.Errorf("no commands found in current folder or in rover home directory")
	}

	var ymlDefinition yamlDefinition

	err = yaml.UnmarshalStrict(commandsFileContent, &ymlDefinition)
	if err != nil {
		return nil, fmt.Errorf("invalid yaml in %s. Internal Error:%s", commandsFilePath, err.Error())
	}

	err = validateCustomCommands(ymlDefinition.Commands)
	if err != nil {
		return nil, err
	}

	for commandName, c := range ymlDefinition.Commands {
		commandList := make([]Command, 1)
		commandList[0] = ymlDefinition.Commands[commandName]

		params := ""
		for i, v := range c.Parameters {
			if i == 0 {
				params += v.Name
			} else {
				params += fmt.Sprintf(",%s", v.Name)
			}
		}

		command := Action{
			Commands: commandList,
			ActionBase: landingzone.ActionBase{
				Name:        commandName,
				Type:        landingzone.CustomCommand,
				Description: fmt.Sprintf("Perform %s with %s parameters", c.ExecutableName, params),
			},
		}

		commands = append(commands, command)
	}

	for groupName, g := range ymlDefinition.Groups {
		commandList := make([]Command, len(ymlDefinition.Groups[groupName]))
		for i, commandName := range ymlDefinition.Groups[groupName] {
			commandList[i] = Command{
				ExecutableName: "rover",
				SubCommand:     commandName,
				SetupEnv:       true,
			}
		}

		params := ""
		for i, v := range g {
			if i == 0 {
				params += v
			} else {
				params += fmt.Sprintf(",%s", v)
			}
		}

		group := Action{
			Commands: commandList,
			ActionBase: landingzone.ActionBase{
				Name:        groupName,
				Type:        landingzone.GroupCommand,
				Description: fmt.Sprintf("Perform %s commands sequentially", params),
			},
		}
		err = validateGroups(ymlDefinition.Groups, commands)
		if err != nil {
			return nil, err
		}

		commands = append(commands, group)
	}

	return commands, nil
}

// Execute runs this custom command by running the external executable
func (a Action) Execute(o *landingzone.Options) error {
	console.Successf("Running custom command: %s %s\n", a.GetName(), o.SourcePath)
	args := []string{}

	for _, command := range a.Commands {

		// TODO : check if the init command has been run
		//if command.RequiresInit {
		//}

		if command.Debug {
			args = append(args, "--debug")
		}

		if command.SubCommand != "" {
			args = append(args, command.SubCommand)
		}

		if command.Flags != "" {
			args = append(args, command.Flags)
		}

		for _, parameter := range command.Parameters {
			templateName := fmt.Sprintf("arguments for action %s", a.GetName())
			argTemplate, err := template.New(templateName).Parse(parameter.Value)
			cobra.CheckErr(err)

			var templateResult bytes.Buffer
			err = argTemplate.Execute(&templateResult, a)
			cobra.CheckErr(err)
			args = append(args, templateResult.String())
		}

		if command.SetupEnv {
			err := o.SetupEnvironment()
			cobra.CheckErr(err)
		}

		// Now ready to actually run it
		cmd := commandpkg.NewCommand(command.ExecutableName, args)
		cmd.Silent = false
		err := cmd.Execute()

		console.Error(cmd.StdErr)
		console.Success(cmd.StdOut)

		// NOTE: When running across multiple levels/stacks
		// We will exit early when we hit first error, this could be improved
		cobra.CheckErr(err)
	}

	return nil
}

func validateCustomCommands(customCommands map[string]Command) error {
	for commandName := range customCommands {
		exists := contains(actions.ActionMap, commandName)

		if exists {
			return fmt.Errorf("invalid custom command (%s). Custom command (%s) cannot be used as it is a builtin command", commandName, commandName)
		}
	}

	return nil
}

func validateGroups(groups map[string][]string, commands []landingzone.Action) error {
	for groupName, group := range groups {
		exists := contains(actions.ActionMap, groupName)
		if exists {
			return fmt.Errorf("invalid group name (%s). (%s) cannot be used as it is a builtin command", groupName, groupName)
		}

		if len(group) == 0 {
			return fmt.Errorf("invalid group (%s). A group must have at least one command", groupName)
		}

		for _, commandName := range group {
			existsBuiltIn := contains(actions.ActionMap, commandName)
			existsCustom := commandsContain(commands, commandName)
			if !existsBuiltIn && !existsCustom {
				return fmt.Errorf("invalid group name (%s). (%s) must be a valid built in command or a custom command", commandName, commandName)
			}
		}
	}

	return nil
}

func commandsContain(commands []landingzone.Action, group string) bool {
	for _, command := range commands {
		if command.GetName() == group {
			return true
		}
	}
	return false
}

func contains(arr map[string]landingzone.Action, str string) bool {
	for _, a := range arr {
		if a.GetName() == str {
			return true
		}
	}
	return false
}
