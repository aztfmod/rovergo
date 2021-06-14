
- `launchpad` - To deploy a launchpad
  - Actions init, plan, apply, destroy, fmt & validate implemented and working
  - Handling of state initialization and upload
  - Handling of destroying a launchpad
- `landingzone` - To deploy a landingzone
  - Actions init, plan, apply, destroy, fmt & validate implemented and working
  - Handling of state and connecting to launchpad
- `landingzone fetch` - Implements the existing `--clone` option
- `cd` - To run actions against single or multiple levels based on symphony config
- `ci <task>` - To run any tasks defined in the ci_tasks directory, which are dynamically discovered
- Supported actions for all `launchpad`, `landingzone`, `cd` are:
  - `init`, `plan`, `apply`, `destroy`, `fmt`, `validate`, 
- Shape of commands & sub-commands and CLI structure defined.
- Handling of locating remote state from level and CAF environment .
- Config file support, currently `.rover.yaml` is used and looked for in $HOME or cwd
- Calling Azure APIs to make calls e.g. query resources with ARG, get storage account, upload blobs, get KV secrets
- Interaction with Azure CLI to obtain subscription and current identity details
- Fundamentals: Goreleaser, GitHub Actions (for CI and release), linting, makefile

### [📝 See the wiki for further details](https://github.com/aztfmod/rovergo/wiki)

# Major Outstanding Work

In very rough order of priority

- Testing with managed identity (system and user)
- Testing with state and deployment in different subscriptions
- Testing using other Azure clouds other than public
- User impersonation
- Terraform Cloud support
- Landingzone list
