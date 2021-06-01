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
			Use: taskConfig.Name,
			Run: func(cmd *cobra.Command, args []string) {
				symphonyConfigFileName, _ := cmd.Parent().PersistentFlags().GetString("symphony-config")
				symphonyConfig, err := symphony.NewSymphonyConfig(symphonyConfigFileName)
				cobra.CheckErr(err)

				debug, _ := cmd.Flags().GetBool("debug")

				if debug {
					symphonyConfig.OutputDebug(symphonyConfigFileName)
					taskConfig.OutputDebug()
				}

				console.Infof("Running ci task %s\n", taskConfig.Name)

			},
		}

		cmd.AddCommand(ciTaskCommand)
	}

}
