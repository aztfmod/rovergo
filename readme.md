# Rover v2

A project to undertake re-writing the [Rover tool](https://github.com/aztfmod/rover) in Go

Uses [Cobra](https://github.com/spf13/cobra) to provide the framework for a robust and familiar CLI tool and [Viper](https://github.com/spf13/viper) for configuration

This new version is not backward compatible with the previous version, and has many large scale changes mainly to the CLI and how it is run, the actual core driving of terraform for CAF landing zone management remains the same

Highlights:
- New "modern" CLI, which feels much like other CLI tools (terraform, helm etc)
- Easier to understand command structure, much easier to run, less error prone
- Standalone binary (for Linux, Mac & Windows) which can be run from anywhere on the system, without the need for hardcoded paths
- Better logging and output
- Use either via the standalone commands `launchpad` or `landingzone` or in a CI/CD pipeline in conjunction with a symphony config file

# Current Status

### [👷‍♂️ Project board](https://github.com/orgs/aztfmod/projects/28?card_filter_query=label%3Arover-go)
### ☢ This is under heavy development, expect breaking changes almost daily 🔥

## Implemented and working

- `launchpad` - To manage a launchpad
  - Actions init, plan, apply, destroy, fmt & validate implemented and working
  - Handling of state initialization and upload
  - Handling of destroying a launchpad
- `landingzone` - To manage a landingzone
  - Actions init, plan, apply, destroy, fmt & validate implemented and working
  - Handling of state and connecting to launchpad
- `cd` - To run actions against single or multiple levels based on symphony config
- `ci <task>` - To run any tasks defined in the ci_tasks directory, which are dynamically discovered
- Supported actions for all `launchpad`, `landingzone`, `cd` commands are:
  - `init`, `plan`, `apply`, `destroy`, `fmt`, `validate`, 
- Handling of locating remote state from level and CAF environment
- Config file support, currently `.rover.yaml` is used and looked for in $HOME or cwd
- Calling Azure APIs to make calls e.g. query resources with ARG, get storage account, upload blobs, get KV secrets
- Interaction with Azure CLI to obtain subscription and current identity details
- Fundamentals: Goreleaser, GitHub Actions (for CI and release), linting, makefile
- `landingzone fetch` - Implements the existing `--clone` option

### [📝 See the wiki for further details](https://github.com/aztfmod/rovergo/wiki)

# Major Outstanding Work

In rough order of priority

- Investigate running with managed identity (system and user)
- Integrating terratest, so that `rover ci test` can be run
- Investigate running with state and deployment in different subscriptions
- Investigate running using other Azure clouds other than public
- User impersonation
- Terraform Cloud support
- Landingzone list
