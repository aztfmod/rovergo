//
// Rover - Azure
// * Ben C, May 2021
//

package azure

import (
	"strings"

	"github.com/spf13/cobra"
)

// CloudNameToTerraform maps cloud names from Azure CLI and SDK names to Terraform names
// * Because no one can agree what to call these!
// * The Azure CLI, the Azure SDK and Terraform all use different names 💩
func CloudNameToTerraform(name string) string {
	switch name {
	case "AzureCloud":
		return "public"
	case "AzurePublicCloud":
		return "public"
	case "AzureChinaCloud":
		return "china"
	case "AzureUSGovernment":
		return "usgovernment"
	case "AzureGermanCloud":
		return "german"
	default:
		return "public"
	}
}

// ParseResourceID into subscription, resource group and name
func ParseResourceID(resourceID string) (subID string, resGrp string, resName string) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		cobra.CheckErr("Supplied resource ID has insufficient segments")
	}

	return parts[2], parts[4], parts[8]
}
