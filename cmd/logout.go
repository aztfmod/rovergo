//
// Rover - Logout command
// * logs out from terraform/az.
// * Greg O, May 2021
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the Azure account.",
	Long:  `Log out from the Azure account.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logout called")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
