package cmd

import (
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cloneCmd represents the clone command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the Azure account",
	Long:  `Authenticate with an Azure account, either the locally logged in user (from Azure CLI), a service principal or managed service identity`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupAzureEnvironment()
		_, err := IsAuthenticated()
		if err != nil {
			cobra.CheckErr(color.RedString("%v", err))
		}
		saveFlags()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	azureEnvDefault, set := os.LookupEnv("ARM_ENVIRONMENT")
	if !set {
		azureEnvDefault = "public"
	}
	useMsiDefaultString, set := os.LookupEnv("ARM_USE_MSI")
	if !set {
		useMsiDefaultString = "false"
	}
	useMsiDefault, _ := strconv.ParseBool(useMsiDefaultString)
	loginCmd.Flags().StringP("subscription-id", "s", os.Getenv("ARM_SUBSCRIPTION_ID"), "The subscription ID which should be used")
	loginCmd.Flags().StringP("client-id", "u", os.Getenv("ARM_CLIENT_ID"), "The client ID which should be used")
	loginCmd.Flags().StringP("client-secret", "p", os.Getenv("ARM_CLIENT_SECRET"), "The client secret which should be used. For use when authenticating as a service principal using a client secret")
	loginCmd.Flags().StringP("tenant-id", "t", os.Getenv("ARM_TENANT_ID"), "The tenant ID which should be used")
	loginCmd.Flags().String("environment", azureEnvDefault, "The cloud environment which should be used. Possible values are public, usgovernment, german, and china")
	loginCmd.Flags().String("msi-endpoint", os.Getenv("ARM_MSI_ENDPOINT"), "The path to a custom endpoint for managed service identity - in most circumstances this should be detected automatically")
	loginCmd.Flags().String("client-certificate-password", os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"), "Certificate password, if client-certificate-path is set")
	loginCmd.Flags().String("client-certificate-path", os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"), "The path to the client certificate associated with the service principal for use when authenticating as a service principal using a client certificate")
	loginCmd.Flags().Bool("use-msi", useMsiDefault, "Allowed managed service identity be used for authentication")
}

func SetupAzureEnvironment() {
	os.Setenv("ARM_SUBSCRIPTION_ID", viper.GetString("subscription-id"))
	os.Setenv("ARM_CLIENT_ID", viper.GetString("client-id"))
	os.Setenv("ARM_TENANT_ID", viper.GetString("tenant-id"))
	os.Setenv("ARM_ENVIRONMENT", viper.GetString("environment"))
	os.Setenv("ARM_CLIENT_CERTIFICATE_PATH", viper.GetString("client-certificate-path"))
	os.Setenv("ARM_CLIENT_CERTIFICATE_PASSWORD", viper.GetString("client-certificate-password"))
	os.Setenv("ARM_CLIENT_SECRET", viper.GetString("client-secret"))
	os.Setenv("ARM_USE_MSI", viper.GetString("use-msi"))
	os.Setenv("ARM_MSI_ENDPOINT", viper.GetString("msi-endpoint"))
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
