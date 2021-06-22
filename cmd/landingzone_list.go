//
// Rover - Landing zone list command
// * This carries out listing of deployed landing zones by querying storage
// * TODO: Stub
// * Ben C, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var lzListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployed landingzones",

	Run: func(cmd *cobra.Command, args []string) {
		level, _ := cmd.Flags().GetString("level")
		cafEnv, _ := cmd.Flags().GetString("environment")
		workspace, _ := cmd.Flags().GetString("workspace")
		stateSub, _ := cmd.Flags().GetString("state-sub")

		if stateSub == "" {
			sub, err := azure.GetSubscription()
			cobra.CheckErr(err)
			stateSub = sub.ID
		}

		console.Infof("Fetching list of all deployed landing zones for level '%s' in env '%s'\n", level, cafEnv)
		launchPadStorageID, _ := azure.FindStorageAccount(level, cafEnv, stateSub)
		if launchPadStorageID == "" {
			console.Error("Unable to locate the launchapd storage account for this level and env")
			cobra.CheckErr("Leaving now...")
		}

		console.Successf("Located launchpad storage account:\n%s\n", launchPadStorageID)
		console.Warning("Landing zones:")
		blobs, err := azure.ListBlobs(launchPadStorageID, workspace)
		cobra.CheckErr(err)
		for _, blob := range blobs {
			name := blob.Name[:len(blob.Name)-8]
			console.Warningf(" - %s\t%dkb\t%v\n", name, *blob.Properties.ContentLength/1024, blob.Properties.LastModified)
		}
	},
}

func init() {
	lzListCmd.Flags().StringP("level", "l", "level1", "CAF level name")
	lzListCmd.Flags().StringP("environment", "e", "sandpit", "Name of CAF environment")
	lzListCmd.Flags().StringP("workspace", "w", "tfstate", "Name of workspace")
	lzListCmd.Flags().String("state-sub", "", "Azure subscription ID where state is held")
	lzListCmd.Flags().SortFlags = true

	landingzoneCmd.AddCommand(lzListCmd)
}
