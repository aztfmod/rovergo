//
// Rover - Azure
// * Ben C, May 2021
//

package azure

import (
	"strings"

	"github.com/spf13/cobra"
)

type WellKnownCloud struct {
	Name          string
	TerraformName string
	KeyvaultDNS   string
	StorageDNS    string
}

var wellKnownClouds []WellKnownCloud

func init() {
	azure := WellKnownCloud{"AzureCloud", "public", "vault.azure.net", "core.windows.net"}
	azurePublic := WellKnownCloud{"AzurePublicCloud", "public", "vault.azure.net", "core.windows.net"}
	china := WellKnownCloud{"AzureChinaCloud", "china", "vault.azure.cn", "core.chinacloudapi.cn"}
	germany := WellKnownCloud{"AzureGermanCloud", "german", "vault.microsoftazure.de", "core.cloudapi.de"}
	gov := WellKnownCloud{"AzureUSGovernment", "usgovernment", "vault.usgovcloudapi.net", "core.usgovcloudapi.net"}

	wellKnownClouds = []WellKnownCloud{azure, azurePublic, china, germany, gov}
}

func CloudNameToTerraform(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.TerraformName
		}

	}

	return "public"
}

func KeyvaultDNSForSubscription() string {
	sub := GetSubscription()
	return KeyvaultDNSForCloud(sub.EnvironmentName)
}

func KeyvaultDNSForCloud(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.KeyvaultDNS
		}

	}

	return ""
}

func StorageDNSForSubscription() string {
	sub := GetSubscription()
	return StorageDNSForCloud(sub.EnvironmentName)
}

func StorageDNSForCloud(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.StorageDNS
		}

	}

	return ""
}

// ParseResourceID into subscription, resource group and name
func ParseResourceID(resourceID string) (subID string, resGrp string, resName string) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		cobra.CheckErr("Supplied resource ID has insufficient segments")
	}

	return parts[2], parts[4], parts[8]
}
