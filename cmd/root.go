//
// Rover - Entry point and root command
// * Handles global flags, initialization, config files and when user runs rover without sub command
//

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aztfmod/rover/pkg/builtin/actions"
	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/custom"
	"github.com/aztfmod/rover/pkg/landingzone"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/aztfmod/rover/pkg/symphony"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/aztfmod/rover/pkg/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rover",
	Short: "Rover is a tool to assist the deployment of the Azure CAF Terraform landingzones",
	Long: `Azure CAF rover is a command line tool in charge of the deployment of the landing zones in your 
Azure environment.
It acts as a toolchain development environment to avoid impacting the local machine but more importantly 
to make sure that all contributors in the GitOps teams are using a consistent set of tools and version.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		console.DebugEnabled, _ = cmd.Flags().GetBool("debug")
	},
}

var helpCmd = &cobra.Command{
	Use:         "help",
	Short:       "Help about any command",
	Long:        "Help about any command",
	Annotations: map[string]string{"cmd_group_annotation": landingzone.BuiltinCommand},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = version.Value
	cobra.CheckErr(rootCmd.Execute())
}

func GetVersion() string {
	return rootCmd.Version
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "log extra debug information, may contain secrets")

	rootCmd.SetHelpCommand(helpCmd)

	command.ValidateDependencies()

	// Ensure rover home exists and create the default contents
	_, err := rover.HomeDirectory()
	if err != nil {
		console.Errorf("Problem with rover home directory: %s\n", err)
		os.Exit(1)
	}

	// Find and load in custom commands
	err = custom.InitializeCustomCommandsAndGroups()
	if err != nil {
		console.Warningf("No custom command or group found in the current directory or rover home directory\n")
	} else {
		console.Infof("Custom commands and groups loaded from %s\n", utils.GetCustomCommandsAndGroupsYamlFilePath())
	}

	BuildSubCommandsFromActionMap()

	fn := func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Short)
		fmt.Println(cmd.Long)
		fmt.Println()

		fmt.Println("Usage:")
		fmt.Printf("  %s [command]\n\n", cmd.Use)

		fmt.Println("Flags:")
		fmt.Println(cmd.Flags().FlagUsages())

		usage := helpMessageByGroups(cmd)
		fmt.Println(usage)
		fmt.Println()

		fmt.Println("Use \"rover [command] --help\" for more information about a command.")
		fmt.Println()
	}

	rootCmd.SetHelpFunc(fn)
}

func BuildSubCommandsFromActionMap() {
	// Dynamically build sub-commands from list of actions
	for key, action := range actions.ActionMap {
		actionSubCmd := &cobra.Command{
			Use:         key,
			Short:       action.GetDescription(),
			Annotations: map[string]string{"cmd_group_annotation": action.GetType()},
			PreRun: func(cmd *cobra.Command, args []string) {
			},
			Run: func(cmd *cobra.Command, args []string) {
				// NOTE: We CAN NOT use the action variable from the loop above as it's not bound at runtime
				// Dynamically building our commands has some limitations, instead we need to use the cmd name & the map
				action = actions.ActionMap[cmd.Name()]

				configFile, _ := cmd.Flags().GetString("config-file")
				configPath, _ := cmd.Flags().GetString("config-dir")

				// Handle the user trying to use both configPath and configFile or neither!
				if configPath == "" && configFile == "" {
					_ = cmd.Help()
					os.Exit(0)
				}
				if configPath != "" && configFile != "" {
					cobra.CheckErr("--config-file and --config-dir options must not be combined, specify only one")
				}

				var optionsList []landingzone.Options
				// Handle symphony mode where config file and level is passed, this will return optionsList with MANY items
				if configFile != "" {
					// Depending on if we're running single or multi-level this will return one or many options
					utils.SymphonyYamlFilePath = configFile
					optionsList = symphony.BuildOptions(cmd)
				}

				// Handle CLI or standalone mode, this will return optionsList with a single item
				if configPath != "" {
					optionsList = landingzone.BuildOptions(cmd)
				}

				for _, options := range optionsList {
					// Now start the action execution...
					// If an error occurs, depend on downstream code to log messages
					console.Infof("Executing action %s for %s\n", action.GetName(), options.StateName)
					err := action.Execute(&options)
					if err != nil {
						cobra.CheckErr(err)
					}
				}

				console.Success("Rover has finished")
				os.Exit(0)
			},
		}

		actionSubCmd.Flags().StringP("source", "s", "", "Path to source of landingzone")
		actionSubCmd.Flags().StringP("config-file", "c", "", "Configuration file, you must supply this or config-dir")
		actionSubCmd.Flags().StringP("config-dir", "v", "", "Configuration directory, you must supply this or config-file")
		actionSubCmd.Flags().StringP("environment", "e", "", "Name of CAF environment")
		actionSubCmd.Flags().StringP("workspace", "w", "", "Name of workspace")
		actionSubCmd.Flags().StringP("statename", "n", "", "Name for state and plan files")
		actionSubCmd.Flags().String("state-sub", "", "Azure subscription ID where state is held")
		actionSubCmd.Flags().String("target-sub", "", "Azure subscription ID to operate on")
		actionSubCmd.Flags().Bool("launchpad", false, "Run in launchpad mode, i.e. level 0")
		actionSubCmd.Flags().StringP("level", "l", "", "CAF landingzone level name, default is all levels")
		actionSubCmd.Flags().BoolP("dry-run", "d", false, "Execute a dry run where no actions will be executed")
		actionSubCmd.Flags().StringP("stack", "t", "", "CAF landingzone level stack name")
		actionSubCmd.Flags().StringP("test-source", "", "", "Path to source of tests")
		actionSubCmd.Flags().SortFlags = true

		// Stuff it under the parent root command
		rootCmd.AddCommand(actionSubCmd)
	}
}

func helpMessageByGroups(cmd *cobra.Command) string {
	groups := map[string][]string{}
	for _, c := range cmd.Commands() {
		var groupName string
		v, ok := c.Annotations["cmd_group_annotation"]
		if !ok {
			groupName = "Other Commands:"
		} else {
			groupName = fmt.Sprintf("%s Commands:", v)
		}

		groupCmds := groups[groupName]
		groupCmds = append(groupCmds, fmt.Sprintf("  %-16s%s", c.Name(), c.Short))
		sort.Strings(groupCmds)

		groups[groupName] = groupCmds
	}

	groupNames := []string{}
	for k := range groups {
		groupNames = append(groupNames, k)
	}
	sort.Strings(groupNames)

	buf := bytes.Buffer{}
	for _, groupName := range groupNames {
		commands := groups[groupName]

		buf.WriteString(fmt.Sprintf("%s\n", groupName))

		for _, cmd := range commands {
			buf.WriteString(fmt.Sprintf("%s\n", cmd))
		}
		buf.WriteString("\n")
	}
	return strings.TrimSpace(buf.String())
}
