//
// Rover - Azure
// * Ben C, May 2021
//

package azure

import (
	"strings"

	"github.com/spf13/cobra"
)

type wellKnownCloud struct {
	Name            string
	TerraformName   string
	KeyvaultDNS     string
	StorageEndpoint string
}

var wellKnownClouds []wellKnownCloud

func init() {
	azure := wellKnownCloud{"AzureCloud", "public", "vault.azure.net", "core.windows.net"}
	azurePublic := wellKnownCloud{"AzurePublicCloud", "public", "vault.azure.net", "core.windows.net"}
	china := wellKnownCloud{"AzureChinaCloud", "china", "vault.azure.cn", "core.chinacloudapi.cn"}
	germany := wellKnownCloud{"AzureGermanCloud", "german", "vault.microsoftazure.de", "core.cloudapi.de"}
	gov := wellKnownCloud{"AzureUSGovernment", "usgovernment", "vault.usgovcloudapi.net", "core.usgovcloudapi.net"}

	wellKnownClouds = []wellKnownCloud{azure, azurePublic, china, germany, gov}
}

func CloudNameToTerraform(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.TerraformName
		}

	}

	return "public"
}

func KeyvaultEndpointForSubscription() string {
	sub := GetSubscription()
	return KeyvaultEndpointForCloud(sub.EnvironmentName)
}

func KeyvaultEndpointForCloud(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.KeyvaultDNS
		}

	}

	return ""
}

func StorageEndpointForSubscription() string {
	sub := GetSubscription()
	return StorageEndpointForCloud(sub.EnvironmentName)
}

func StorageEndpointForCloud(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.StorageEndpoint
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
