package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cloneCmd represents the clone command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the azure account",
	Long:  `Login into the azure account`,
	Run: func(cmd *cobra.Command, args []string) {
		setArmEnv()
		_, err := IsAuthenticated()
		if err != nil {
			log.Fatal(err)
		} else {
			saveFlags()
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
	loginCmd.Flags().StringP("subscription_id", "s", os.Getenv("ARM_SUBSCRIPTION_ID"), "The Subscription ID which should be used")
	loginCmd.Flags().StringP("client_id", "u", os.Getenv("ARM_CLIENT_ID"), "The Client ID which should be used")
	loginCmd.Flags().StringP("client_secret", "p", os.Getenv("ARM_CLIENT_SECRET"), "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret")
	loginCmd.Flags().String("tenant_id", os.Getenv("ARM_TENANT_ID"), "The Tenant ID which should be used")
	loginCmd.Flags().String("environment", env_default, "The Cloud Environment which should be used. Possible values are public, usgovernment, german, and china")
	loginCmd.Flags().String("msi_endpoint", os.Getenv("ARM_MSI_ENDPOINT"), "the path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically")
	loginCmd.Flags().String("client_certificate_password", os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"), "certificate password")
	loginCmd.Flags().String("client_certificate_path", os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"), "the path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate")
	loginCmd.Flags().Bool("use_msi", use_msi_default, "allowed Managed Service Identity be used for Authentication")
}

func setArmEnv() {
	os.Setenv("ARM_SUBSCRIPTION_ID", viper.GetString("subscription_id"))
	os.Setenv("ARM_CLIENT_ID", viper.GetString("client_id"))
	os.Setenv("ARM_TENANT_ID", viper.GetString("tenant_id"))
	os.Setenv("ARM_ENVIRONMENT", viper.GetString("environment"))
	os.Setenv("ARM_CLIENT_CERTIFICATE_PATH", viper.GetString("client_certificate_path"))
	os.Setenv("ARM_CLIENT_CERTIFICATE_PASSWORD", viper.GetString("client_certificate_password"))
	os.Setenv("ARM_CLIENT_SECRET", viper.GetString("client_secret"))
	os.Setenv("ARM_USE_MSI", viper.GetString("use_msi"))
	os.Setenv("ARM_MSI_ENDPOINT", viper.GetString("msi_endpoint"))
}

func saveFlags() {
	viper.WriteConfig()
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
