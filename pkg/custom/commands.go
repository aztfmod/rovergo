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
