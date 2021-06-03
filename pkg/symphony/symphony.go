package symphony

import (
	"fmt"
	"os"

	"github.com/aztfmod/rover/pkg/console"
	"gopkg.in/yaml.v2"
)

type Config struct {
	FileName string
	Content  struct {
		Environment  string `yaml:"environment,omitempty"`
		Repositories []struct {
			Name   string `yaml:"name,omitempty"`
			URI    string `yaml:"uri,omitempty"`
			Branch string `yaml:"branch,omitempty"`
		}
		Levels []struct {
			Level     string `yaml:"level,omitempty"`
			Type      string `yaml:"type,omitempty"`
			Launchpad bool   `yaml:"launchpad,omitempty"`
			Stacks    []struct {
				Stack             string `yaml:"stack,omitempty"`
				LandingZonePath   string `yaml:"landingZonePath,omitempty"`
				ConfigurationPath string `yaml:"configurationPath,omitempty"`
				TfState           string `yaml:"tfState,omitempty"`
			}
		}
	}
}

func NewSymphonyConfig(symphonyConfigFileName string) (*Config, error) {
	sc := new(Config)
	sc.FileName = symphonyConfigFileName

	buf, _ := os.ReadFile(symphonyConfigFileName)
	err := yaml.Unmarshal(buf, &sc.Content)

	return sc, err
}

func (sc *Config) OutputDebug() {
	fmt.Println()

	console.Debugf("Verbose output of %s\n", sc.FileName)
	console.Debugf(" - Environment: %s\n", sc.Content.Environment)
	console.Debugf(" - Number of repositories: %d\n", len(sc.Content.Repositories))
	console.Debugf(" - Number of levels: %d\n", len(sc.Content.Levels))
}
