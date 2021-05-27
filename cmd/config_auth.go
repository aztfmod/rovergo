//
// Rover - Config/auth sub-command
// * Evaluates & stores authentication properties such as client-id, client-secret
// * Greg O, May 2021
//

package cmd

import (
	"os"
	"strconv"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/terraform"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultEnvironment = "public"
const defaultUseMsi = "false"

// cfgAuthCmd represents the auth command
var cfgAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Evaluate authentication configuration parameters.",
	Long: `Evaluate and store authentication configuration parameters,
such as client-id, client-secret and so on. Stores the configuration
into a local file (by default ./.rover.yaml). With the --clear flag,
clear that configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Handle when user passes `--clear` and wipe saved values
		if doClear, _ := cmd.Flags().GetBool("clear"); doClear {
			console.Warning("Clearing stored credentials and service principal details")
			for k := range viper.GetStringMap("auth") {
				viper.Set("auth."+k, "")
			}
			// Reset some values to defaults
			defaultMsi := false
			defaultMsi, _ = strconv.ParseBool(defaultUseMsi)
			viper.Set("auth.use-msi", defaultMsi)
			viper.Set("auth.environment", defaultEnvironment)
			saveFlags()
			return
		}

		_, err := terraform.Authenticate()

		if err != nil {
			cobra.CheckErr(err)
		}

		saveFlags()

		for k, v := range viper.GetStringMapString("auth") {
			console.Debugf("%s = %s\n", k, v)
		}
	},
}

func init() {
	configCmd.AddCommand(cfgAuthCmd)

	azureEnvDefault, set := os.LookupEnv("ARM_ENVIRONMENT")
	if !set {
		azureEnvDefault = defaultEnvironment
	}
	useMsiDefaultString, set := os.LookupEnv("ARM_USE_MSI")
	if !set {
		useMsiDefaultString = defaultUseMsi
	}
	useMsiDefault, _ := strconv.ParseBool(useMsiDefaultString)
	cfgAuthCmd.Flags().StringP("subscription-id", "s", os.Getenv("ARM_SUBSCRIPTION_ID"), "Subscription ID which should be used")
	cfgAuthCmd.Flags().StringP("client-id", "u", os.Getenv("ARM_CLIENT_ID"), "Client ID which should be used")
	cfgAuthCmd.Flags().StringP("client-secret", "p", os.Getenv("ARM_CLIENT_SECRET"), "Client secret which should be used. For use when authenticating as a service principal using a client secret")
	cfgAuthCmd.Flags().StringP("tenant-id", "t", os.Getenv("ARM_TENANT_ID"), "Tenant ID which should be used")
	cfgAuthCmd.Flags().String("environment", azureEnvDefault, "Azure cloud environment which should be used. Possible values are public, usgovernment, german, and china")
	cfgAuthCmd.Flags().String("msi-endpoint", os.Getenv("ARM_MSI_ENDPOINT"), "Path to a custom endpoint for managed service identity - in most circumstances this should be detected automatically")
	cfgAuthCmd.Flags().String("client-cert-password", os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"), "Certificate password, if --client-cert-path is set")
	cfgAuthCmd.Flags().String("client-cert-path", os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"), "Path to the service principal's (spn) client certificate for use when authenticating as an spn using a client certificate")
	cfgAuthCmd.Flags().Bool("use-msi", useMsiDefault, "Try to use managed service identity for authentication")
	cfgAuthCmd.Flags().Bool("clear", false, "Reset and clear any stored credentials")

	// Very important we bind flags to config, and put under the 'auth.' section key
	cfgAuthCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "clear" {
			return
		}
		_ = viper.BindPFlag("auth."+f.Name, f)
	})
}

func saveFlags() {
	err := viper.WriteConfig()
	cobra.CheckErr(err)
}
