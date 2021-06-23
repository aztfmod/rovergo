package custom

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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
	SetupEnv    bool `yaml:"setupEnv"`
}

// This is used to provide the main things you'd want to refer to in a template expression
type argTemplateContext struct {
	Options landingzone.Options
	Action  Action
	Meta    map[string]string
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
func (a Action) Execute(o *landingzone.Options) error {
	console.Successf("Running custom action: %s %s\n", a.Name, o.SourcePath)
	args := []string{}

	if a.command.SetupEnv {
		err := o.SetupEnvironment()
		cobra.CheckErr(err)
	}

	// This allows for golang templated expressions in command arguments
	// e.g. "--foo={{ .Options.SourcePath }}" see https://golang.org/pkg/text/template/
	for _, argDefined := range a.command.Arguments {
		templateName := fmt.Sprintf("arguments for action %s", a.Name)
		argTemplate, err := template.New(templateName).Parse(argDefined)
		cobra.CheckErr(err)

		roverDir, err := rover.HomeDirectory()
		cobra.CheckErr(err)

		// Build conext to execute the template with
		templateContext := argTemplateContext{
			Options: *o,
			Action:  a,
			Meta: map[string]string{
				"RoverHome": roverDir,
			},
		}

		var templateResult bytes.Buffer
		err = argTemplate.Execute(&templateResult, templateContext)
		cobra.CheckErr(err)
		args = append(args, templateResult.String())
	}

	// Now ready to actually run it
	cmd := command.NewCommand(a.command.Executable, args)
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
	roverHomeDir, err := rover.HomeDirectory()
	if err != nil {
		return nil, err
	}
	custActionsPath := filepath.Join(roverHomeDir, actionsFilename)
	// _, err = os.Stat(custActionsPath)
	// // If doesn't exist then place our default YAML file in .rover
	// if err != nil {
	// 	fileErr := ioutil.WriteFile(custActionsPath, []byte(defaultFileContent), 0777)
	// 	if fileErr != nil {
	// 		return nil, fileErr
	// 	}
	// }

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
