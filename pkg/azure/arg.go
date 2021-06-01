//
// Rover - Azure resource graph
// * Supports querying Azure resource graph (ARG)
// * Ben C, May 2021
//

package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2019-04-01/resourcegraph"
	"github.com/spf13/cobra"
)

// RunQuery against the Azure resource graph
func RunQuery(query string, subID string) interface{} {
	argClient := resourcegraph.New()
	argClient.Authorizer = GetAuthorizer()

	requestOpts := resourcegraph.QueryRequestOptions{
		ResultFormat: "objectArray",
	}

	request := resourcegraph.QueryRequest{
		Subscriptions: &[]string{subID},
		Query:         &query,
		Options:       &requestOpts,
	}

	// Run the query and get the results
	var results, err = argClient.Resources(context.Background(), request)
	cobra.CheckErr(err)

	return results.Data
}
