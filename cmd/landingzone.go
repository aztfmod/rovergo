package cmd

import (
	"github.com/spf13/cobra"
)

// landingzoneCmd represents the landingzone command
var landingzoneCmd = &cobra.Command{
	Use:     "landingzone",
	Aliases: []string{"lz"},
	Short:   "Manage and deploy landing zones",
	Long:    `This command allows you to deploy, update and destory CAF landing zones`,
}

func init() {
	rootCmd.AddCommand(landingzoneCmd)

}
