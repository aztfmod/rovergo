//
// Rover - test command
// * Experimental code work in progress used to test and call other code paths & functions
// * Ben C, May 2021
//

package cmd

import (
	"github.com/aztfmod/rover/pkg/azure"
	"github.com/spf13/cobra"
)

// TODO: Experimental code work in progress
// ****************************************

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		azure.UploadFileToBlob(
			"/subscriptions/52512f28-c6ed-403e-9569-82a9fb9fec91/resourceGroups/odog-rg-launchpad-level0/providers/Microsoft.Storage/storageAccounts/odogstlevel0",
			"tfstate",
			"dave.tfstate",
			"/home/ben/tfstates/level0/tfstate/mystate.tfstate")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
