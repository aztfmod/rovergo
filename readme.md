# Rover Go

A project to undertake re-writing the [Rover tool](https://github.com/aztfmod/rover) in Go

Uses [Cobra](https://github.com/spf13/cobra) to provide the framework for a robust and familiar CLI tool and [Viper](https://github.com/spf13/viper) for configuration

# Current Status

### ☢ This is under heavy development, expect braking changes almost daily 🔥

## Implemented 

- Shape of commands & sub-commands and CLI structure
- `launchpad fetch` - Implements the existing `--clone` option
- `launchpad run` - To deploy a launchpad
  - Actions init, plan, deploy & destroy implemented and working
  - Handling of state initialization and upload
  - Handling of locating remote state from level and CAF environment
- `cd run` - To run actions against multiple levels based on symphony config (only level0 currently)
- `ci <task>` - To run any tasks defined in the ci_tasks directory, which are dynamically discovered
- Config file support, currently `.rover.yaml` is used and looked for in $HOME or cwd
- Calling Azure APIs to make calls e.g. query resources with ARG, get storage account, upload blobs
- Interaction with Azure CLI to obtain subscription and current identity details
- Goreleaser, GitHub Actions, linting, makefile

### [📝 See the wiki for further details](https://github.com/aztfmod/rovergo/wiki)

```text
$ rover
Azure CAF rover is a command line tool in charge of the deployment of the landing zones in your 
Azure environment.
It acts as a toolchain development environment to avoid impacting the local machine but more importantly 
to make sure that all contributors in the GitOps teams are using a consistent set of tools and version.

Usage:
  rover [command]

Available Commands:
  cd          Manage CD operations.
  ci          Manage CI operations.
  help        Help about any command
  landingzone Manage and deploy landing zones
  launchpad   Manage and deploy launchpad, i.e. landing zone level0.
  terraform   Manage terraform operations.
  workspace   Manage workspace operations.

Flags:
      --config string   config file (default is ./.rover.yaml)
      --debug           log extra debug information, may contain secrets
  -h, --help            help for rover
  -v, --version         version for rover

Use "rover [command] --help" for more information about a command.
```

The flags on the launchpad / landingzone are as follows:

```text
Run actions to deploy, update or remove launchpads

Usage:
  rover launchpad run [flags]

Flags:
  -a, --action string        Action to run, one of [init | plan | deploy | destroy] (default "init")
  -c, --config-path string   Configuration vars directory (required)
  -e, --environment string   Name of CAF environment (default "sandpit")
  -h, --help                 help for run
  -s, --source string        Path to source of landingzone (required)
      --state-sub string     Azure subscription ID where state is held
  -n, --statename string     Name for state and plan files, defaults to landingzone name
      --target-sub string    Azure subscription ID to operate on
  -w, --workspace string     Name of workspace (default "tfstate")
```
# Major Outstanding Work

In very rough order of priority

- Rest of launchpad and landing zone operations - priority #1 ☺
- Other remote state cases, e.g. login_as_launchpad
- Testing with managed identity (system and user)
- User impersonation
- CI operations
- CD operations
- Terraform Cloud support
