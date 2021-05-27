# Rover Go

A project to undertake re-writing the [Rover tool](https://github.com/aztfmod/rover) in Go

Uses [Cobra](https://github.com/spf13/cobra) to provide the framework for a robust and familiar CLI tool and [Viper](https://github.com/spf13/viper) for configuration

# Current Status

### ☢ This should be considered spike / POC / investigation grade code 🔥

## Implemented 

- Shape of commands & sub-commands and CLI structure
- `rover launchpad fetch` - Implements the existing `--clone` option
- `rover config auth` - To configure which Azure credentials to use
- Config file support, currently `.rover.yaml` is used and looked for in $HOME or cwd
- Minimal implementation of Terraform init / plan / apply
- Authentication into Terraform and Azure API SDK
  - Service Principal with secret or cert - if configured
  - Managed Identity - if configured
  - Azure CLI - default if above not configured
- Calling Azure API to make calls e.g. get storage account
- Goreleaser, GitHub Actions, linting, makefile
 
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
  config      Access to configuration related sub-commands, such as 'auth'.
  help        Help about any command
  landingzone Manage and deploy landing zones
  launchpad   Manage and deploy launchpad, i.e. landing zone level0.
  logout      Log out from the Azure account.
  terraform   Manage terraform operations.
  workspace   Manage workspace operations.

Flags:
      --config string   config file (default is ./.rover.yaml)
      --debug           log extra debug information, may contain secrets
  -h, --help            help for rover
  -v, --version         version for rover

Use "rover [command] --help" for more information about a command.
```

# Major Outstanding Work

In very rough order of priority

- All launchpad and landing zone operations - priority #1 ☺
- All handling remote state
- User impersonation
- CI operations
- CD operations
- Terraform Cloud support
