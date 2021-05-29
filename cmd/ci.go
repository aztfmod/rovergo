//
// Rover - Top level ci (continuous integration) command
// * Doesn't do anything, all work is done by sub-commands
// * Greg O, May 2021
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/symphony"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Manage CI operations.",
	Long:  `Manage CI operations.`,
	Run: func(cmd *cobra.Command, args []string) {

		symphonyConfigFileName, _ := cmd.Flags().GetString("symphony-config")
		debug, _ := cmd.Flags().GetBool("debug")

		symphonyConfig, err := symphony.NewSymphonyConfig(symphonyConfigFileName)
		cobra.CheckErr(err)

		if debug {
			symphonyConfig.OutputDebug(symphonyConfigFileName)
		}

		run(symphonyConfigFileName)
	},
}

func run(symphonyConfigFileName string) {
	fmt.Println()

	console.Infof("Running CI command, config: %s\n", symphonyConfigFileName)
}

func init() {
	ciCmd.Flags().StringP("symphony-config", "c", "", "Path/filename of symphony.yml")
	ciCmd.Flags().SetNormalizeFunc(aliasNormalizeFunc)

	ciCmd.Flags().BoolP("verbose", "v", false, "Output symphony.yml to console")

	err := cobra.MarkFlagRequired(ciCmd.Flags(), "symphony-config")
	cobra.CheckErr(err)

	rootCmd.AddCommand(ciCmd)

	addCITasks(ciCmd)
}

func aliasNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "sc":
		name = "symphony-config"
	}
	return pflag.NormalizedName(name)
}

func addCITasks(cmd *cobra.Command) {

	debug, _ := cmd.Flags().GetBool("debug")

	ciTaskConfigFileNames, err := getCITaskConfigFilenames("./ci_tasks")
	if err != nil {
		return
	}

	for _, filename := range ciTaskConfigFileNames {

		taskConfig, err := symphony.NewTaskConfig(filepath.Join("./ci_tasks", filename))
		cobra.CheckErr(err)

		if debug {
			taskConfig.OutputDebug(filename)
		}

		var ciTaskCommand = &cobra.Command{
			Use: taskConfig.Name,
			Run: func(cmd *cobra.Command, args []string) {
				console.Infof("Running ci task %s", taskConfig.Name)
				if taskConfig.SubCommand != "" {
					console.Infof(" - with sub-command %s", taskConfig.SubCommand)
				}
			},
		}

		cmd.AddCommand(ciTaskCommand)
	}

}

func getCITaskConfigFilenames(directoryName string) ([]string, error) {
	var files []string

	f, err := os.Open(directoryName)
	if err != nil {
		return files, err
	}

	fileInfo, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil

}
