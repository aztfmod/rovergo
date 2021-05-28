package azure

import (
	"encoding/json"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"
)

// Account holds details fetched from `az account show` command
type Account struct {
	EnvironmentName  string
	TenantID         string
	SubscriptionName string `json:"name"`
	SubscriptionID   string `json:"id"`
	User             map[string]string
}

// GetAccount gets the current logged in details from the Azure CLI
// Will fail and exit if they aren't found
func GetAccount() Account {
	err := command.CheckCommand("az")
	cobra.CheckErr(err)

	accountRes, err := command.QuickRun("az", "account", "show", "-o=json")
	cobra.CheckErr(err)

	account := &Account{}
	err = json.Unmarshal([]byte(accountRes), account)
	cobra.CheckErr(err)

	console.Successf("Azure account details obtained for user: %s\n", account.User["name"])
	console.Successf("Azure subscription is: %s (%s)\n", account.SubscriptionName, account.SubscriptionID)
	return *account
}
