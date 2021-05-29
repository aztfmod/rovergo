package symphony

import (
	"fmt"
	"os"

	"github.com/aztfmod/rover/pkg/console"
	"gopkg.in/yaml.v2"
)

type TaskConfig struct {
	Name           string `yaml:"environment,omitempty"`
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
	buf, _ := os.ReadFile(taskConfigFileName)
	err := yaml.Unmarshal(buf, p)

	return p, err
}

func (tc *TaskConfig) OutputDebug(taskConfigFileName string) {
	fmt.Println()

	console.Debugf("Verbose output of ci task configuration %s\n", taskConfigFileName)
	console.Debugf(" - Task name: %s\n", tc.Name)
	console.Debugf(" - Executable name: %d\n", tc.ExecutableName)
}
