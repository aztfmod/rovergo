package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rover",
	Short: "Rover is a tool to assist the deployment of the Azure CAF Terraform landingzones",
	Long: `Azure CAF rover is a command line tool in charge of the deployment of the landing zones in your 
Azure environment.
It acts as a toolchain development environment to avoid impacting the local machine but more importantly 
to make sure that all contributors in the GitOps teams are using a consistent set of tools and version.`,
	PersistentPreRun: verify_version,
}

func verify_version(cmd *cobra.Command, args []string) {
	//TODO: Verify that the current version installed is the latest version kind of what pip does
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	cobra.CheckErr(rootCmd.Execute())

	debug, _ = rootCmd.Flags().GetBool("debug")
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rover.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "emit debug information, may contain secrets")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory and CWD with name ".rover" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rover")

		// Config defaults
		viper.SetDefault("tempDir", "/tmp")
		viper.SetDefault("terraform.install", true)
		viper.SetDefault("terraform.install-path", "./bin")
		viper.SetDefault("state.storage-account", "")
		viper.SetDefault("state.container", "")
		viper.SetDefault("state.resource-group", "")
		viper.SetDefault("state.access-key", "")
	}

	viper.SetEnvPrefix("rover")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		color.Green("Using config file: %s", viper.ConfigFileUsed())
	} else {
		// Fall back to creating empty config file
		fileName := "./.rover.yaml"
		_, err := os.Create(fileName)
		cobra.CheckErr(err)
		color.Yellow("Config file not found, creating new file %s with defaults", fileName)
	}
}
