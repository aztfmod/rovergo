package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// landingzoneCmd represents the landingzone command
var landingzoneCmd = &cobra.Command{
	Use:     "landingzone",
	Aliases: []string{"lz"},
	Short:   "Manage landing zones",
	Long:    `Blah blah `,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nARGS: %+v\n", args)
		fmt.Println("landingzone called")
	},
}

func init() {
	rootCmd.AddCommand(landingzoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// landingzoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// landingzoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
