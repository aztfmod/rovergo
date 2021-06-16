//
// Rover - Azure CLI
// * Interactions with the Azure CLI
// * Ben C, May 2021
//

package azure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
)

// User holds details of the signed in user, might be a managed identity
// populated with values from 'az account show'
type User struct {
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
	User            User
}

// Identity holds an Azure AD identity; user, SP or MSI
type Identity struct {
	UserPrincipalName string
	ObjectType        string
	ObjectID          string
	Mail              string
	MailNickname      string
	DisplayName       string
}

type UserAssignedIdentityIDs struct {
	ClientID    string `json:"clientID,omitempty"`
	PrincipalID string `json:"principalID,omitempty"`
}

// VMIdentity is the output of 'az vm identity show'
type VMIdentity struct {
	PrincipalID            string                     `json:"principalID,omitempty"`
	TenantID               string                     `json:"tenantID,omitempty"`
	IdentityType           string                     `json:"type,omitempty"`
	UserAssignedIdentities map[string]json.RawMessage `json:"userAssignedIdentities,omitempty"`
}

// BasicIdentity - can be either User or ServicePrincipal
type BasicIdentity struct {
	DisplayName string
	ObjectID    string
	ObjectType  string
	ClientID    string
}

type VMIdentities struct {
	IDList []BasicIdentity
}

// GetSubscription gets the current logged in details from the Azure CLI
// Will fail and exit if they aren't found
func GetSubscription() Subscription {
	err := command.CheckCommand("az")
	cobra.CheckErr(err)

	cmdRes, err := command.QuickRun("az", "account", "show", "-o=json")
	cobra.CheckErr(err)

	sub := &Subscription{}
	err = json.Unmarshal([]byte(cmdRes), sub)
	cobra.CheckErr(err)

	console.Successf("Azure subscription is: %s (%s)\n", sub.Name, sub.ID)
	return *sub
}

// GetIdentity gets the current logged in user from the Azure CLI
// Will fail and exit if they aren't found
func GetIdentity() BasicIdentity {
	err := command.CheckCommand("az")
	cobra.CheckErr(err)

	cmdRes, err := command.QuickRun("az", "ad", "signed-in-user", "show", "-o=json")
	cobra.CheckErr(err)

	ident := &Identity{}
	err = json.Unmarshal([]byte(cmdRes), ident)
	cobra.CheckErr(err)

	basicIdent := BasicIdentity{
		DisplayName: ident.DisplayName,
		ObjectID:    ident.ObjectID,
		ObjectType:  ident.ObjectType,
	}
	console.Successf("Signed in indentity is '%s' (%s)\n", ident.UserPrincipalName, ident.ObjectType)
	return basicIdent
}

// GetVMIdentities will get the MI details of an Azure VM, both system assigned and user assigned
func GetVMIdentities(resourceGroupName string, vmName string) VMIdentities {
	err := command.CheckCommand("az")
	cobra.CheckErr(err)

	cmdRes, err := command.QuickRun(
		"az",
		"vm",
		"identity",
		"show",
		fmt.Sprintf("--resource-group=%s", resourceGroupName),
		fmt.Sprintf("--name=%s", vmName),
		"-o=json")
	cobra.CheckErr(err)

	vmident := &VMIdentity{}
	err = json.Unmarshal([]byte(cmdRes), vmident)
	cobra.CheckErr(err)

	var identities VMIdentities
	if strings.Contains(vmident.IdentityType, "SystemAssigned") {
		identities.IDList = append(identities.IDList, BasicIdentity{
			DisplayName: "SystemAssigned",
			ObjectType:  "servicePrincipal",
			ObjectID:    vmident.PrincipalID,
		})
	}

	if vmident.UserAssignedIdentities != nil {

		for _, uai := range vmident.UserAssignedIdentities {

			ids := &UserAssignedIdentityIDs{}
			err = json.Unmarshal([]byte(uai), ids)
			cobra.CheckErr(err)

			identities.IDList = append(identities.IDList, BasicIdentity{
				DisplayName: "UserAssigned",
				ObjectType:  "servicePrincipal",
				ObjectID:    ids.PrincipalID,
				ClientID:    ids.ClientID,
			})

		}
	}

	return identities
}
