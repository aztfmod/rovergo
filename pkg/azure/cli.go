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
	"github.com/spf13/cobra"
)

// Subscription holds details fetched from `az account show` command
type Subscription struct {
	EnvironmentName string
	TenantID        string
	Name            string
	ID              string
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
func GetIdentity() Identity {
	err := command.CheckCommand("az")
	cobra.CheckErr(err)

	cmdRes, err := command.QuickRun("az", "ad", "signed-in-user", "show", "-o=json")
	cobra.CheckErr(err)

	ident := &Identity{}
	err = json.Unmarshal([]byte(cmdRes), ident)
	cobra.CheckErr(err)

	console.Successf("Signed in indentity is '%s' (%s)\n", ident.UserPrincipalName, ident.ObjectType)
	return *ident
}
