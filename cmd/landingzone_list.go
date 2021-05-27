package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var lzListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all landingzones",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("landingzones list called")
	},
}

func init() {
	landingzoneCmd.AddCommand(lzListCmd)
}
