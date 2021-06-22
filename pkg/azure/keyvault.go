//
// Rover - Azure keyvault
// * For working with azure keyvault getting secrets
// * Greg O, June 2021
//

package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type SdkKVClient struct {
	AudienceHostName string
	VaultName        string
	Client           keyvault.BaseClient
}

func NewKVClient(audienceHostName string, vaultName string) (*SdkKVClient, error) {

	kvc := new(SdkKVClient)
	kvc.AudienceHostName = audienceHostName
	kvc.VaultName = vaultName

	kvc.Client = keyvault.New()
	authorizer, err := auth.NewAuthorizerFromCLIWithResource(kvc.Audience())
	if err != nil {
		return nil, err
	}
	kvc.Client.Authorizer = authorizer

	return kvc, nil
}

func (kvc *SdkKVClient) Audience() string {
	return fmt.Sprintf("https://%s", kvc.AudienceHostName)
}

func (kvc *SdkKVClient) VaultBaseURL() string {
	return fmt.Sprintf("https://%s.%s", kvc.VaultName, kvc.AudienceHostName)
}

func (kvc *SdkKVClient) GetSecret(secretName string) (string, error) {

	result, err := kvc.Client.GetSecret(context.Background(), kvc.VaultBaseURL(), secretName, "")
	if err != nil {
		return "", err
	}

	secretValue := result.Value

	return *secretValue, nil
}

// FindKeyVault returns resource id for a CAF tagged keyvault
func FindKeyVault(level string, environment string, subID string) (string, error) {
	query := fmt.Sprintf(`Resources 
		| where type == 'microsoft.keyvault/vaults' 
		| where tags.level == '%s'
		| where tags.environment == '%s'
		| limit 1
		| project id`, level, environment)

	queryResults, err := RunQuery(query, subID)
	if err != nil {
		return "", err
	}

	resSlice, ok := queryResults.([]interface{})
	if !ok {
		return "", errors.New("findKeyVault: error asserting query results")
	}

	if len(resSlice) <= 0 {
		return "", errors.New("no keyvault found")
	}

	resMap, ok := resSlice[0].(map[string]interface{})
	if !ok {
		return "", errors.New("findKeyVault: Failed to parse query results")
	}

	return resMap["id"].(string), nil
}
