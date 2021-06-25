//
// Rover - Azure
// * Ben C, May 2021
//

package azure

import (
	"fmt"
	"strings"
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

func KeyvaultEndpointForSubscription() (string, error) {
	sub, err := GetSubscription()
	if err != nil {
		return "", err
	}
	return KeyvaultEndpointForCloud(sub.EnvironmentName), nil
}

func KeyvaultEndpointForCloud(name string) string {

	for _, cloud := range wellKnownClouds {

		if cloud.Name == name {
			return cloud.KeyvaultDNS
		}

	}

	return ""
}

func StorageEndpointForSubscription() (string, error) {
	sub, err := GetSubscription()
	if err != nil {
		return "", err
	}
	return StorageEndpointForCloud(sub.EnvironmentName), nil
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
// TODO: There's a function the Azure SDK cli package that we can replace this with I think
func ParseResourceID(resourceID string) (subID string, resGrp string, resName string, err error) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		return "", "", "", fmt.Errorf("Supplied resource ID has insufficient segments")
	}

	return parts[2], parts[4], parts[8], err
}
