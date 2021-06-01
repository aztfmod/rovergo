//
// Rover - Azure auth
// * Helper functions for signing into Azure to use the Go SDK
// * Ben C, May 2021
//

package azure

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
)

// GetAuthorizer used for Azure SDK access logs into Azure
// This should always be called before any Azure SDK calls
func GetAuthorizer() autorest.Authorizer {
	// We defer everything to the Azure CLI
	azureAuthorizer, err := auth.NewAuthorizerFromCLI()
	cobra.CheckErr(err)

	return azureAuthorizer
}
