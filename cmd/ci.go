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

	addCITasks(ciCmd)
}

func addCITasks(cmd *cobra.Command) {

	directoryName := "./ci_tasks"
	pTaskConfigs, err := symphony.NewTaskConfigs(directoryName)
	cobra.CheckErr(err)

	for _, filename := range pTaskConfigs.EnumerateFilenames() {

		taskConfig, err := symphony.NewTaskConfig(filepath.Join(directoryName, filename))
		cobra.CheckErr(err)

		var ciTaskCommand = &cobra.Command{
			Use: taskConfig.Name,
			Run: func(cmd *cobra.Command, args []string) {
				symphonyConfigFileName, _ := cmd.Flags().GetString("symphony-config")
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

		ciTaskCommand.Flags().StringP("symphony-config", "c", "", "Path/filename of symphony.yml")
		err = cobra.MarkFlagRequired(ciTaskCommand.Flags(), "symphony-config")
		cobra.CheckErr(err)

		cmd.AddCommand(ciTaskCommand)
	}

}
