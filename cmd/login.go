package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/spf13/cobra"
)

var subscription_id, client_id, client_secret, tenant_id, environment, msi_endpoint, client_certificate_password, client_certificate_path string
var use_msi bool

// cloneCmd represents the clone command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the azure account",
	Long:  `Login into the azure account`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := login()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	env_default, set := os.LookupEnv("ARM_ENVIRONMENT")
	if !set {
		env_default = "public"
	}
	use_msi_default_string, set := os.LookupEnv("ARM_USE_MSI")
	if !set {
		use_msi_default_string = "false"
	}
	use_msi_default, _ := strconv.ParseBool(use_msi_default_string)
	loginCmd.Flags().StringVar(&subscription_id, "subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID"), "The Subscription ID which should be used")
	loginCmd.Flags().StringVar(&client_id, "client_id", os.Getenv("ARM_CLIENT_ID"), "The Client ID which should be used")
	loginCmd.Flags().StringVar(&client_id, "client_secret", os.Getenv("ARM_CLIENT_SECRET"), "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret")
	loginCmd.Flags().StringVar(&tenant_id, "tenant_id", os.Getenv("ARM_TENANT_ID"), "The Tenant ID which should be used")
	loginCmd.Flags().StringVar(&environment, "environment", env_default, "The Cloud Environment which should be used. Possible values are public, usgovernment, german, and china")
	loginCmd.Flags().StringVar(&msi_endpoint, "msi_endpoint", os.Getenv("ARM_MSI_ENDPOINT"), "the path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically")
	loginCmd.Flags().StringVar(&client_certificate_password, "client_certificate_password", os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"), "certificate password")
	loginCmd.Flags().StringVar(&client_certificate_path, "client_certificate_path", os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"), "the path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate")
	loginCmd.Flags().BoolVar(&use_msi, "use_msi", use_msi_default, "allowed Managed Service Identity be used for Authentication")
}

func login() (*authentication.Config, error) {
	os.Setenv("ARM_SUBSCRIPTION_ID", subscription_id)
	os.Setenv("ARM_CLIENT_ID", client_id)
	os.Setenv("ARM_TENANT_ID", tenant_id)
	os.Setenv("ARM_ENVIRONMENT", environment)
	os.Setenv("ARM_CLIENT_CERTIFICATE_PATH", client_certificate_path)
	os.Setenv("ARM_CLIENT_CERTIFICATE_PASSWORD", client_certificate_password)
	os.Setenv("ARM_CLIENT_SECRET", client_secret)
	os.Setenv("ARM_USE_MSI", strconv.FormatBool(use_msi))
	os.Setenv("ARM_MSI_ENDPOINT", msi_endpoint)
	return IsAuthenticated()
}

func IsAuthenticated() (*authentication.Config, error) {
	builder := &authentication.Builder{
		TenantOnly:               false,
		SupportsAuxiliaryTenants: false,
		AuxiliaryTenantIDs:       nil,
		SupportsAzureCliToken:    true,
		SupportsClientCertAuth:   true,
		SupportsClientSecretAuth: true,
	}
	return builder.Build()
}
