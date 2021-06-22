//
// Rover - Azure resource graph
// * Supports querying Azure resource graph (ARG)
// * Ben C, May 2021
//

package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2019-04-01/resourcegraph"
)

// RunQuery against the Azure resource graph
func RunQuery(query string, subID string) (interface{}, error) {
	argClient := resourcegraph.New()
	authorizer, err := GetAuthorizer()
	if err != nil {
		return nil, err
	}
	argClient.Authorizer = authorizer

	requestOpts := resourcegraph.QueryRequestOptions{
		ResultFormat: "objectArray",
	}

	request := resourcegraph.QueryRequest{
		Subscriptions: &[]string{subID},
		Query:         &query,
		Options:       &requestOpts,
	}

	// Run the query and get the results
	results, err := argClient.Resources(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return results.Data, nil
}
