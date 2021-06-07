//
// Rover - Landing zone list command
// * This carries out listing of deployed landing zones by querying storage
// * TODO: Stub
// * Ben C, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var lzListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployed landingzones",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("landingzones list is not implemented")
	},
}

func init() {
	landingzoneCmd.AddCommand(lzListCmd)
}
