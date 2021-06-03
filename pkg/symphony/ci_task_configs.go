package symphony

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aztfmod/rover/pkg/console"
	"gopkg.in/yaml.v2"
)

type TaskConfigs struct {
	Filenames []string
}

type TaskConfig struct {
	FileName string
	Content  struct {
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
}

func NewTaskConfig(taskConfigFileName string) (*TaskConfig, error) {
	tc := new(TaskConfig)
	tc.FileName = taskConfigFileName

	buf, _ := os.ReadFile(taskConfigFileName)
	err := yaml.Unmarshal(buf, &tc.Content)

	tc.Content.Name = strings.ToLower(tc.Content.Name)

	return tc, err
}

func (tc *TaskConfig) OutputDebug() {
	fmt.Println()

	console.Debugf("Verbose output of ci task configuration, file name: %s\n", tc.FileName)
	console.Debugf(" - Task name: %s\n", tc.Content.Name)
	console.Debugf(" - Executable name: %s\n", tc.Content.ExecutableName)
	if tc.Content.SubCommand != "" {
		console.Debugf(" - Sub-command name: %s\n", tc.Content.SubCommand)
	}
}

func FindTaskConfig(directoryName string, taskName string) (*TaskConfig, error) {

	pTaskConfigs, err := NewTaskConfigs(directoryName)
	if err != nil {
		return nil, err
	}

	var foundTaskConfig = new(TaskConfig)
	for _, filename := range pTaskConfigs.EnumerateFilenames() {

		taskConfig, err := NewTaskConfig(filepath.Join(directoryName, filename))
		if err != nil {
			return nil, err
		}

		if strings.EqualFold(taskConfig.Content.Name, taskName) {
			foundTaskConfig = taskConfig
			break
		}
	}

	return foundTaskConfig, nil
}

func NewTaskConfigs(directoryName string) (*TaskConfigs, error) {

	tcs := new(TaskConfigs)

	f, err := os.Open(directoryName)
	if err != nil {
		return nil, err
	}

	fileInfo, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		tcs.Filenames = append(tcs.Filenames, file.Name())
	}

	return tcs, nil

}

func (tcs *TaskConfigs) EnumerateFilenames() []string {
	return tcs.Filenames
}
