//
// Rover - Azure CLI
// * Interactions with the Azure CLI
// * Ben C, May 2021
//

package azure

import (
	"encoding/json"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
)

// AccountUser holds details of the signed in user, might be a managed identity
// populated with values from 'az account show'
type AccountUser struct {
	AssignedIdentityInfo string `json:"assignedIdentityInfo,omitempty"`
	Name                 string `json:"name,omitempty"`
	Usertype             string `json:"type,omitempty"`
}

// Subscription holds details fetched from `az account show` command
type Subscription struct {
	EnvironmentName string
	TenantID        string
	Name            string
	ID              string
	User            AccountUser
}

// Identity holds an Azure AD identity; user, SP or MSI
type signedInUserIdentity struct {
	UserPrincipalName string
	ObjectType        string
	ObjectID          string
	Mail              string
	MailNickname      string
	DisplayName       string
}

// VMIdentity is the output of 'az vm identity show'
type VMIdentity struct {
	PrincipalID            string                     `json:"principalID,omitempty"`
	TenantID               string                     `json:"tenantID,omitempty"`
	IdentityType           string                     `json:"type,omitempty"`
	UserAssignedIdentities map[string]json.RawMessage `json:"userAssignedIdentities,omitempty"`
}

// Identity - can be either User or ServicePrincipal
type Identity struct {
	DisplayName string
	ObjectID    string
	ObjectType  string
	ClientID    string
}

// GetSubscription gets the current logged in details from the Azure CLI
// Will fail and exit if they aren't found
func GetSubscription() (*Subscription, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	cmdRes, err := command.QuickRun("az", "account", "show", "-o=json")
	if err != nil {
		return nil, err
	}

	sub := &Subscription{}
	err = json.Unmarshal([]byte(cmdRes), sub)
	if err != nil {
		return nil, err
	}

	console.Successf("Azure subscription is: %s (%s)\n", sub.Name, sub.ID)
	return sub, nil
}

// GetSignedInIdentity gets the current logged in user from the Azure CLI
// Will fail and exit if they aren't found
// Will Fail if az is authenticated with a service principal. Use the GetSignedInIdentityServicePrincipal function instead
func GetSignedInIdentity() (*Identity, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	cmdRes, err := command.QuickRun("az", "ad", "signed-in-user", "show", "-o=json")
	if err != nil {
		return nil, err
	}

	ident := &signedInUserIdentity{}
	err = json.Unmarshal([]byte(cmdRes), ident)
	if err != nil {
		return nil, err
	}

	basicIdent := &Identity{
		DisplayName: ident.DisplayName,
		ObjectID:    ident.ObjectID,
		ObjectType:  ident.ObjectType,
	}
	return basicIdent, nil
}

// GetSignedInIdentity gets the current logged in service principal from the Azure CLI
// note az ad signed-in-user show does not work for sp's. see https://github.com/Azure/azure-cli/issues/10439
func GetSignedInIdentityServicePrincipal() (*Identity, error) {
	//TODO: az account show
	//TODO: az sp show --id <id from az account show user>
	//TODO: return Identity struct with DisplayName, ClientId and ObjectId of SP.
	return nil, nil
}
