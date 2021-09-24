package custom

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/builtin/actions"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/rover"
	"gopkg.in/yaml.v2"
)

const commandsFileName = "commands.yml"

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
	Parameters     []struct {
		Name   string `yaml:"name"`
		Value  bool   `yaml:"value"`
		Prefix string `yaml:"prefix"`
	} `yaml:"parameters"`
}

// LoadCustomCommandsAndGroups is called by root cmd during init
// It finds all the custom action definitions and returns them to be plugged into the CLI
func LoadCustomCommandsAndGroups() (commands []landingzone.Action, err error) {
	// Getting the current working directory
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// e.g. /github/example_project/commands.yml
	commandsFilePath := filepath.Join(currentWorkingDirectory, commandsFileName)

	var fileInfo os.FileInfo

	// Checks if the commands file exists in the current working directory
	if fileInfo, err = os.Stat(commandsFilePath); os.IsNotExist(err) {

		// If the file does not exist, get the rover home directory
		roverHomeDir, err := rover.HomeDirectory()
		if err != nil {
			return nil, err
		}
		// e.g. ~/.rover/commands.yml
		commandsFilePath = filepath.Join(roverHomeDir, commandsFileName)

		// Checks if the commands file exists in the current working directory
		if fileInfo, err = os.Stat(commandsFilePath); os.IsNotExist(err) {
			// If the file does not exist, return an empty list of commands
			// and Not Exists Error
			return nil, os.ErrNotExist
		}
	}

	// if the file exists, but empty, return an empty list of commands
	if fileInfo.Size() == 0 {
		return nil, nil
	}

	commandsFileContent, err := ioutil.ReadFile(commandsFilePath)
	if err != nil {
		return nil, err
	}

	var ymlDefinition yamlDefinition

	err = yaml.UnmarshalStrict(commandsFileContent, &ymlDefinition)
	if err != nil {
		return nil, err
	}

	err = validateCustomCommands(ymlDefinition.Commands)
	if err != nil {
		return nil, err
	}
	err = validateGroups(ymlDefinition.Groups)
	if err != nil {
		return nil, err
	}

	return
}

func validateCustomCommands(customCommands map[string]Command) error {
	for commandName := range customCommands {
		exists := contains(actions.ActionMap, commandName)

		if exists {
			return fmt.Errorf("custom command name (%s) cannot be the same as a builtin command", commandName)
		}
	}

	return nil
}

func validateGroups(groups map[string][]string) error {
	for groupName, group := range groups {
		exists := contains(actions.ActionMap, groupName)

		if exists {
			return fmt.Errorf("group name (%s) cannot be the same as a builtin command", groupName)
		}

		for _, commandName := range group {
			exists := contains(actions.ActionMap, commandName)

			if !exists {
				return fmt.Errorf("group command name (%s) must be exist in builtin commands", commandName)
			}
		}
	}

	return nil
}

func contains(arr map[string]landingzone.Action, str string) bool {
	for _, a := range arr {
		if a.GetName() == str {
			return true
		}
	}
	return false
}
