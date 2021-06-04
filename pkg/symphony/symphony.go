package symphony

import (
	"errors"
	"fmt"
	"os"

	"github.com/aztfmod/rover/pkg/console"
	"gopkg.in/yaml.v2"
)

type Config struct {
	FileName string
	Content  struct {
		Version         int    `yaml:"symphonyVersion,omitempty"`
		Environment     string `yaml:"environment,omitempty"`
		LandingZonePath string `yaml:"landingZonePath,omitempty"`
		Workspace       string
		Repositories    []struct {
			Name   string `yaml:"name,omitempty"`
			URI    string `yaml:"uri,omitempty"`
			Branch string `yaml:"branch,omitempty"`
		}
		Levels []Level
	}
}

type Level struct {
	Name      string `yaml:"level,omitempty"`
	Type      string `yaml:"type,omitempty"`
	Launchpad bool   `yaml:"launchpad,omitempty"`
	Stacks    []Stack
}

type Stack struct {
	Name              string `yaml:"stack,omitempty"`
	LandingZonePath   string `yaml:"landingZonePath,omitempty"`
	ConfigurationPath string `yaml:"configurationPath,omitempty"`
	TfState           string `yaml:"tfState,omitempty"`
}

func NewSymphonyConfig(symphonyConfigFileName string) (*Config, error) {
	sc := new(Config)
	buf, _ := os.ReadFile(symphonyConfigFileName)
	err := yaml.Unmarshal(buf, &sc.Content)
	if sc.Content.Version != 2 {
		return nil, errors.New("bad symphony version number, this version of rover requires version 2")
	}

	return sc, err
}

func (sc *Config) OutputDebug() {
	fmt.Println()

	console.Debugf("Verbose output of %s\n", sc.FileName)
	console.Debugf(" - Environment: %s\n", sc.Content.Environment)
	console.Debugf(" - Number of repositories: %d\n", len(sc.Content.Repositories))
	console.Debugf(" - Number of levels: %d\n", len(sc.Content.Levels))
}
