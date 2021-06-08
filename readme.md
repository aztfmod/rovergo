# Rover v2

A project to undertake re-writing the [Rover tool](https://github.com/aztfmod/rover) in Go

Uses [Cobra](https://github.com/spf13/cobra) to provide the framework for a robust and familiar CLI tool and [Viper](https://github.com/spf13/viper) for configuration

# Current Status

### [👷‍♂️ Project board](https://github.com/orgs/aztfmod/projects/28?card_filter_query=label%3Arover-go)
### ☢ This is under heavy development, expect braking changes almost daily 🔥

## Implemented 

- Shape of commands & sub-commands and CLI structure
- `launchpad` - To deploy a launchpad
  - Actions init, plan, apply, destroy, fmt & validate implemented and working
  - Handling of state initialization and upload
  - Handling of locating remote state from level and CAF environment
  - Handling of destorying a launchpad
- `landingzone` - To deploy a landingzone
  - Actions init, plan, apply, destroy, fmt & validate implemented and working
  - Handling of state and connecting to launchpad
- `landingzone fetch` - Implements the existing `--clone` option
- `cd` - To run actions against multiple levels based on symphony config
- `ci <task>` - To run any tasks defined in the ci_tasks directory, which are dynamically discovered
- Config file support, currently `.rover.yaml` is used and looked for in $HOME or cwd
- Calling Azure APIs to make calls e.g. query resources with ARG, get storage account, upload blobs, get KV secrets
- Interaction with Azure CLI to obtain subscription and current identity details
- Goreleaser, GitHub Actions, linting, makefile

### [📝 See the wiki for further details](https://github.com/aztfmod/rovergo/wiki)

# Major Outstanding Work

In very rough order of priority

- Testing with managed identity (system and user)
- Testing with state and deployment in different subscriptions
- User impersonation
- Terraform Cloud support
- Landingzone list
