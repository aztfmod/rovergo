package symphony

import (
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/spf13/cobra"
)

func (c Config) RunCITask(ciTasksDirectoryName string, subCommandName string, targetLevel string, debug bool) {

	taskConfig, err := FindTaskConfig(ciTasksDirectoryName, subCommandName)
	cobra.CheckErr(err)

	console.Debugf("Running executable %s, sub-command %s, level %s\n\n", taskConfig.Content.ExecutableName, taskConfig.Content.SubCommand, targetLevel)

	for _, level := range c.Content.Levels {

		if targetLevel == "all" || targetLevel == level.Name {

			for _, stack := range level.Stacks {

				console.Debugf("Running ci task %s in environment %s, level %s, stack %s\n",
					subCommandName,
					c.Content.Environment,
					level.Name,
					stack.Name)

				opt := landingzone.Options{
					LaunchPadMode:  level.Launchpad,
					CafEnvironment: c.Content.Environment,
					Workspace:      c.Content.Workspace,
					Level:          level.Name,
				}
				opt.SetConfigPath(stack.ConfigurationPath)
				opt.SetSourcePath(stack.LandingZonePath)

				if strings.ToLower(taskConfig.Content.ExecutableName) == "terraform" {

					action, err := landingzone.NewAction(taskConfig.Content.SubCommand)
					cobra.CheckErr(err)

					stateName := stack.TfState
					// IMPORTANT: We use the stack name as the default name if tfState key is not supplied
					if stateName == "" {
						stateName = stack.Name
					}
					opt.StateName = stateName

					opt.Execute(action)

				} else {

					// non-terraform commands

					args := []string{}
					for _, arg := range taskConfig.Content.Parameters {
						args = append(args, arg.Prefix+arg.Name)
						args = append(args, arg.Value)
					}

					cmd := command.NewCommand(
						taskConfig.Content.ExecutableName,
						args,
					)

					if debug {
						cmd.DryRun = true
						cmd.Silent = false
					}

					err := cmd.Execute()
					cobra.CheckErr(err)

				}

			}
		}
	}

}
