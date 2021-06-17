package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest/azure/cli"
	"github.com/aztfmod/rover/pkg/console"
)

const graphAPIEndpoint = "https://graph.microsoft.com"

type graphResult struct {
	Value []struct {
		AppDisplayName string
		ID             string
		AppID          string
	}
}

func GetServicePrincipalIdentity(clientID string) (*Identity, error) {
	token, err := cli.GetTokenFromCLI(graphAPIEndpoint)
	if err != nil {
		return nil, err
	}
	console.Debugf("Obtained a token for %s expires on %s\n", graphAPIEndpoint, token.ExpiresOn)

	// This query filter was obtained from the `az ad sp show` command
	filter := fmt.Sprintf(`servicePrincipalNames/any(c:appId eq '%s')`, clientID)
	filter = url.QueryEscape(filter)
	urlString := fmt.Sprintf(`%s/beta/servicePrincipals?$filter=%s&$top=1`, graphAPIEndpoint, filter)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	console.Debugf("Making API call to %s\n", urlString)

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Graph API error %s", resp.Status)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	console.Debugf("API call returned %s and %d bytes\n", resp.Status, len(bodyBytes))
	data := graphResult{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Value) < 1 {
		return nil, fmt.Errorf("No service principal found with client-id %s", clientID)
	}

	return &Identity{
		DisplayName: data.Value[0].AppDisplayName,
		ObjectID:    data.Value[0].ID,
		ClientID:    data.Value[0].AppID,
		ObjectType:  "servicePrincipal",
	}, nil
}
