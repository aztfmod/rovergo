package custom

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/builtin_actions"
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

	return
}

