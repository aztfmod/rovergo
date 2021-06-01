//
// Rover - Azure cloud
// * Ben C, May 2021
//

package azure

// CloudNameToTerraform maps cloud names from Azure CLI and SDK names to Terraform names
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
