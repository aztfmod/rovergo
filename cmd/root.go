//
// Rover - Entry point and root command
// * Handles global flags, initialization, config files and when user runs rover without sub command
// * Ben C, May 2021
//

package cmd

import (
	"os"

	"github.com/aztfmod/rover/pkg/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rover",
	Short: "Rover is a tool to assist the deployment of the Azure CAF Terraform landingzones",
	Long: `Azure CAF rover is a command line tool in charge of the deployment of the landing zones in your 
Azure environment.
It acts as a toolchain development environment to avoid impacting the local machine but more importantly 
to make sure that all contributors in the GitOps teams are using a consistent set of tools and version.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.DebugEnabled, _ = cmd.Flags().GetBool("debug")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rover.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "log extra debug information, may contain secrets")
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
