//
// Rover - Azure auth
// * Helper functions for signing into Azure to use the Go SDK
// * Ben C, May 2021
//

package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
)

// OwnerRoleDefintionID is fixed GUID for the Owner role in Azure
// See docs https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles
const OwnerRoleDefintionID = "8e3af657-a8ff-443c-a75c-2fe8c4bcb635"

// GetAuthorizer used for Azure SDK access logs into Azure
// This should always be called before any Azure SDK calls
func GetAuthorizer() autorest.Authorizer {
	// We defer everything to the Azure CLI
	azureAuthorizer, err := auth.NewAuthorizerFromCLI()
	cobra.CheckErr(err)

	return azureAuthorizer
}

// CheckIsOwner returns if the given objectId is assigned Owner role on the given subscription
func CheckIsOwner(objectID string, subID string) (bool, error) {
	client := authorization.NewRoleAssignmentsClient(subID)
	client.Authorizer = GetAuthorizer()
	resultPages, err := client.ListForScope(context.Background(), fmt.Sprintf("/subscriptions/%s", subID), fmt.Sprintf("assignedTo('%s')", objectID))
	if err != nil {
		return false, err
	}

	for ; resultPages.NotDone(); err = resultPages.Next() {
		if err != nil {
			return false, err
		}
		for _, roleAssignment := range resultPages.Values() {
			if strings.Contains(*roleAssignment.RoleDefinitionID, OwnerRoleDefintionID) {
				return true, nil
			}
		}
	}

	return false, err
}
