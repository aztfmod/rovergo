//
// Rover - Top level ci (continuous integration) command
// * Doesn't do anything, all work is done by sub-commands, the implementations of which
// *  are methods of symphony.Config
// * Greg O, May 2021
//

package cmd

import (
	"path/filepath"

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

	// BUG: #52 Because this is run inside an init() function it means rover won't run
	// - if ci_tasks is missing, even if the command you are running is NOT `rover ci`
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

				symphonyConfig.RunCITask(directoryName, subCommandName, level, debug)

			},
		}

		cmd.AddCommand(ciTaskCommand)
	}

}
