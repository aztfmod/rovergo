//
// Rover - Azure storage
// * For working with azure storage for accessing & storing TF state
// * TODO: !!! Experimental code work in progress !!!
// * Ben C, May 2021
//

package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/spf13/cobra"
)

func FindStorageAccount(level string, environment string, subscription string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := storage.NewAccountsClient(subscription)
	client.Authorizer = GetAuthorizer()

	accountPages, err := client.List(ctx)
	cobra.CheckErr(err)

	for accountPages.NotDone() {
		for _, account := range accountPages.Values() {
			fmt.Println(*account.Name)
		}
		accountPages.Next()
	}
}
