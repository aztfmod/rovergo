# Rover Go

A project to undertake re-writing the [Rover tool](https://github.com/aztfmod/rover) in Go

Uses [Cobra](https://github.com/spf13/cobra) to provide the framework for a robust and familiar CLI tool and [Viper](https://github.com/spf13/viper) for configuration

# Current Status

### ☢ This should be considered spike / POC / investigation grade code 🔥

## Implemented 

- Shape of commands & sub-commands and CLI structure
- `launchpad fetch` - Implements the existing `--clone` option
- `launchpad run` - To deploy a launchpad
  - Actions init, plan and deploy implemented and working
  - Handling of state initialization and upload
  - Handling of locating remote state from level and CAF environment
- `launchpad ci <task>` - To run any tasks defined in the ci_tasks directory
- Config file support, currently `.rover.yaml` is used and looked for in $HOME or cwd
- Calling Azure APIs to make calls e.g. get storage account, upload blobs
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

# Major Outstanding Work

In very rough order of priority

- Rest of launchpad and landing zone operations - priority #1 ☺
- Other remote state cases, e.g. login_as_launchpad
- Testing with managed identity (system and user)
- User impersonation
- CI operations
- CD operations
- Terraform Cloud support
