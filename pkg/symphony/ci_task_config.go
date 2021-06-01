package symphony

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aztfmod/rover/pkg/console"
	"gopkg.in/yaml.v2"
)

type TaskConfig struct {
	Name           string `yaml:"name,omitempty"`
	ExecutableName string `yaml:"executableName,omitempty"`
	SubCommand     string `yaml:"subCommand,omitempty"`
	Flags          string `yaml:"flags,omitempty"`
	Debug          bool   `yaml:"debug,omitempty"`
	RequiresInit   bool   `yaml:"requiresInit,omitempty"`
	Parameters     []struct {
		Name   string `yaml:"name,omitempty"`
		Value  string `yaml:"value,omitempty"`
		Prefix string `yaml:"prefix,omitempty"`
	}
}

func NewTaskConfig(taskConfigFileName string) (*TaskConfig, error) {
	p := new(TaskConfig)
	reader, _ := os.Open(taskConfigFileName)
	buf, _ := ioutil.ReadAll(reader)
	err := yaml.Unmarshal(buf, p)

	return p, err
}

func (tc *TaskConfig) OutputDebug() {
	fmt.Println()

	console.Debugf("Verbose output of ci task configuration\n")
	console.Debugf(" - Task name: %s\n", tc.Name)
	console.Debugf(" - Executable name: %s\n", tc.ExecutableName)
	if tc.SubCommand != "" {
		console.Debugf(" - Sub-command name: %s\n", tc.SubCommand)
	}
}
