//
// Rover - Top level ci (continuous integration) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"path/filepath"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/symphony"
	"github.com/spf13/cobra"
)

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Manage CI operations.",
	Long:  `Manage CI operations.`,
}

func init() {
	rootCmd.AddCommand(ciCmd)

	ciCmd.PersistentFlags().String("ci-task-dir", "./ci_tasks", "Directory containing the ci task definition files.")
	ciCmd.PersistentFlags().StringP("symphony-config", "c", "./symphony.yaml", "Path/filename of symphony.yaml.")
	ciCmd.PersistentFlags().StringP("level", "l", "all", "Landing zone level to run within.")

	addCITasks(ciCmd)
}

func addCITasks(cmd *cobra.Command) {

	// directoryName := "./ci_tasks"
	directoryName, _ := cmd.PersistentFlags().GetString("ci-task-dir")

	pTaskConfigs, err := symphony.NewTaskConfigs(directoryName)
	cobra.CheckErr(err)

	for _, filename := range pTaskConfigs.EnumerateFilenames() {

		taskConfig, err := symphony.NewTaskConfig(filepath.Join(directoryName, filename))
		cobra.CheckErr(err)

		var ciTaskCommand = &cobra.Command{
			Use: taskConfig.Content.Name,
			Run: func(cmd *cobra.Command, args []string) {

				symphonyConfigFileName, _ := cmd.Parent().PersistentFlags().GetString("symphony-config")
				symphonyConfig, err := symphony.NewSymphonyConfig(symphonyConfigFileName)
				cobra.CheckErr(err)

				directoryName, _ := cmd.Parent().PersistentFlags().GetString("ci-task-dir")

				level, _ := cmd.Parent().PersistentFlags().GetString("level")

				subCommandName := cmd.Use

				debug, _ := rootCmd.PersistentFlags().GetBool("debug")

				if debug {
					symphonyConfig.OutputDebug()
					taskConfig.OutputDebug()
				}

				runCITaskSubCommand(directoryName, subCommandName, symphonyConfig, level, debug)

			},
		}

		cmd.AddCommand(ciTaskCommand)
	}

}

func runCITaskSubCommand(ciTasksDirectoryName string, subCommandName string, symphonyConfig *symphony.Config, targetLevel string, debug bool) {

	console.Debugf("Running ci task %s\n", subCommandName)

	taskConfig, err := symphony.FindTaskConfig(ciTasksDirectoryName, subCommandName)
	cobra.CheckErr(err)

	console.Debugf("Running executable %s, sub-command %s, level %s\n\n", taskConfig.Content.ExecutableName, taskConfig.Content.SubCommand, targetLevel)

	for _, level := range symphonyConfig.Content.Levels {

		if targetLevel == "all" || targetLevel == level.Level {

			for _, stack := range level.Stacks {

				lzp := stack.LandingZonePath
				tfs := stack.TfState

				console.Debugf("Running ci task %s with tfstate %s in folder %s\n", subCommandName, tfs, lzp)
			}
		}
	}
}
