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
	// if viper.GetString("auth.client-id") != "" {
	// 	os.Setenv("AZURE_SUBSCRIPTION_ID", viper.GetString("auth.subscription-id"))
	// 	os.Setenv("AZURE_CLIENT_ID", viper.GetString("auth.client-id"))
	// 	os.Setenv("AZURE_TENANT_ID", viper.GetString("auth.tenant-id"))
	// 	// Terraform and the Go SDK use different values here 🙃
	// 	os.Setenv("AZURE_ENVIRONMENT", "azure"+viper.GetString("auth.environment")+"cloud")
	// 	os.Setenv("AZURE_CERTIFICATE_PATH", viper.GetString("auth.client-cert-path"))
	// 	os.Setenv("AZURE_CERTIFICATE_PASSWORD", viper.GetString("auth.client-cert-password"))
	// 	os.Setenv("AZURE_CLIENT_SECRET", viper.GetString("auth.client-secret"))

	// 	azureAuthorizer, err := auth.NewAuthorizerFromEnvironment()
	// 	cobra.CheckErr(err)

	// 	return azureAuthorizer
	// } else if viper.GetBool("auth.use-msi") {
	// 	msiConfig := auth.NewMSIConfig()
	// 	azureAuthorizer, err := msiConfig.Authorizer()
	// 	cobra.CheckErr(err)

	// 	return azureAuthorizer
	// } else {
	// This is fall through, and will fail if user is not logged in
	azureAuthorizer, err := auth.NewAuthorizerFromCLI()
	cobra.CheckErr(err)

	return azureAuthorizer
	// }
}
