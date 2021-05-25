package cmd

import (
	"os"
	"strconv"

	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/fatih/color"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// cloneCmd represents the clone command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the Azure account",
	Long:  `Authenticate with an Azure account, either the locally logged in user (from Azure CLI), a service principal or managed service identity`,
	Run: func(cmd *cobra.Command, args []string) {
		terraform.SetupAzureEnvironment()
		_, err := isAuthenticated()
		if err != nil {
			cobra.CheckErr(color.RedString("%v", err))
		}
		saveFlags()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	// // for the login command map flags to config file keys
	// viper.BindPFlags(loginCmd.LocalFlags())

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
	loginCmd.Flags().String("client-cert-password", os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"), "Certificate password, if client-certificate-path is set")
	loginCmd.Flags().String("client-cert-path", os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"), "The path to the client certificate associated with the service principal for use when authenticating as a service principal using a client certificate")
	loginCmd.Flags().Bool("use-msi", useMsiDefault, "Allowed managed service identity be used for authentication")

	// Important we bind flags to config, and put under the 'auth.' section key
	loginCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag("auth."+f.Name, f)
	})
}

func saveFlags() {
	viper.WriteConfig()
}

func isAuthenticated() (*authentication.Config, error) {
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
