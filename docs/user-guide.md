# Rover v2 User Guide

This is a guide to the Rover v2 command line and operation.

If you are not familiar with CAF landing zones, please [refer to the concepts section below](#caf-concepts)

## The Rover v2 CLI

```text
Usage:
  rover [command]

Flags:
      --debug     log extra debug information, may contain secrets
  -h, --help      help for rover
  -v, --version   version for rover

Builtin Commands:
  apply           Perform a terraform plan & apply
  destroy         Perform a terraform destroy
  fmt             Perform a terraform format
  help            Help about any command
  init            Perform a terraform init and no other action
  landingzone     Manage and deploy landing zones
  plan            Perform a terraform plan
  validate        Perform a terraform validate
```

The Rover v2 CLI maps actions (see below) to top level sub commands, with the addition of a special landingzone command

### Action Commands

All commands other than `landingzone` are action commands, they all take the same form and have the same switches, as follows

```text
Usage:
  rover <action> [flags]

Flags:
  -v, --config-dir string    Configuration directory, you must supply this or config-file
  -c, --config-file string   Configuration file, you must supply this or config-dir
  -d, --dry-run              Execute a dry run where no actions will be executed
  -e, --environment string   Name of CAF environment, default "sandpit"
  -h, --help                 help for init
      --launchpad            Run in launchpad mode, i.e. level0
  -l, --level string         CAF landingzone level name, default is all levels
  -s, --source string        Path to source of landingzone
      --state-sub string     Azure subscription ID where state is held
  -n, --statename string     Name for state and plan files, default is picked based on source dir name
      --target-sub string    Azure subscription ID to operate on
  -w, --workspace string     Name of workspace

Global Flags:
      --debug                Log extra debug information, may contain secrets
```

### Landingzone Management Commands

```text
Usage:
  rover landingzone [command]

Aliases:
  landingzone, lz

Available Commands:
  fetch       Fetch supporting artifacts such as landingzones from GitHub
  list        List all deployed landingzones
```

## Actions

Rover v2 actions take two forms:

- **Terraform actions**, these map 1:1 to the various terraform commands, with the sole exception of ***apply*** which carries out both a ***plan*** followed by ***apply***
  - init
  - plan
  - apply
  - destroy
  - fmt
  - validate
  - test (integration tests in terratest)

- **Custom actions**, these extend the Rover v2 command set, and allow you to run an external command against the given CAF config, level, source etc. [See below](#custom-actions)

## Running Rover v2

Rover has two main modes of operation, which applies across all actions, these modes are "ad-hoc mode" and "multi-level mode", the switches you supply `--config-dir` and `--config-file` determine which mode is used.

In these examples we assume to have a copy of the CAF landingzones in `./landingzones` and our configuration directory at `./caf-config/` there will be sub-directories under there for our various levels

## Running in "Ad-hoc mode" (single level)

This mode is intended for users not wishing to build a YAML configuration, maybe they want to get started quickly or don't need the multi-level features. In this mode you supply a source directory, a config directory, but also all the other parameters required.

**💬 NOTE. This mode is engaged by the use of the `--config-dir` switch**

Examples:

- Running init action for a launch pad

```bash
rover init --config-dir ./caf-config/level0/launchpad --source ./landingzones --launchpad --level level0
```

- Running apply action for a launch pad

```bash
rover apply --config-dir ./caf-config/level0/launchpad --source ./landingzones --launchpad --level level0
```

- Running destroy action for a landing zone with config held in `level1/myapp`

```bash
rover destroy --config-dir ./caf-config/level1/myapp --source ./landingzones --level level1
```

- Running apply action for a landing zone with config held in `level1/myapp` this time we deploy with the CAF environment set to "prod"

```bash
rover apply --config-dir ./caf-config/level1/myapp --source ./landingzones --level level1 --environment prod
```

**👁‍🗨 Warning.** Despite being optional. It is **STRONGLY** recommended to supply both the `--environment` and `--statename` options when running Rover in ad-hoc mode, this will prevent anything unexpected from happening by Rover picking defaults for these values.

## Running in "Config file mode" (multi level)

This mode is intended for use in CI/CD pipelines when multiple levels are being managed at once, and each of those levels contains multiple "stacks". All settings are held within a YAML configuration file (which is [part of project symphony](https://github.com/aztfmod/symphony)). Rover requires either a single level to be specified, or by default all levels are run, all other settings are obtained from the YAML file.

For an example and reference symphony config file see [examples/ref-app-symphony.yaml](./../examples/ref-app-symphony.yaml)

**💬 NOTE. This mode is engaged by the use of the `--config-file` switch**

Examples:

- Running plan for level 2

```bash
rover plan --config-file ./symphony.yaml --level level2
```

- Running apply for all levels

```bash
rover apply --config-file ./symphony.yaml
```

- Running destroy for all levels

```bash
rover destroy --config-file ./symphony.yaml
```

## Switch Reference

### Shared - Switches

- `--level` Set which level is being operated on
- `--dry-run` Set to perform a dry run and output details of the operation without executing it.

### Ad-hoc Mode - Switches

- `--source` The source landingzone repo location
- `--launchpad` Run deployment in launchpad mode, *only* run this with a valid launchpad config and with level set to "level0"
- `--environment` Set the CAF environment, **defaults to "sandpit"**
- `--statename` Set the state name used for naming tfstate files, **defaults to "caf_solution" or "caf_launchpad"** depending if --launchpad is set
- `--target-sub` Subscription ID to deploy resources into, defaults to current set on Azure CLI
- `--state-sub` Subscription ID where state (i.e. the launchpad) is held, defaults to current set on Azure CLI
- `--workspace` Workspace is used to name the containers used for state, **defaults to "tfstate"**

## Custom Actions

RoverGo has an extensible CLI and command set. A file named `commands.yaml` can be used to extend the available commands in RoverGo with custom commands. This file is can contain custom command definitions and groups of commands definitions. It is loaded as RoverGo starts, and can be edited & amended as required.

Rover will search for commands.yaml in either of the following locations:

- The current directory where RoverGo is invoked.
- The [Rover Home Directory](#rover-home-dir) (~/$HOME/.rover).

An [example commands.yml file](../examples/custom_commands/commands.yml) is provided in this repo.
### Root structure of commands.yaml

```yaml
commands:
# list of commands
groups:
# list of groups of commands
```

### Custom Commands Reference

Each top level key in the file is used as the name of a new custom command. In the following example a new command called `finder` is introduced.

```yaml
# This is provided as an example
finder:
  description: "find custom command short description"
  executableName: "find"
  subCommand: "fmt"
  flags: "-no-color -recursive -check -diff"
  debug: false
  requiresInit: false
  parameters:
    - name: list
      value: true
      prefix: "-"
    - name: write
      value: false
      prefix: "-"
```

Each custom command definition supports the following options:

- `description` - short description of the custom command
- `executableName` - The name of executable or command to run, must be on the system path or fully qualified
- `subCommand` - The sub command to run, e.g. `apply`, `test` or `plan`
- `flags` - The flags to pass to the executable, e.g. `-no-color -recursive -check -diff`
- `debug` - A boolean flag to enable debug output, defaults to false
- `requiresInit` - A boolean flag to indicate if Rover needs to be initialised before running the command, defaults to false
- `parameters` - A list of parameters to pass to the executable, e.g. `-list -write`

The parameters field can be static strings but also supports [Go templating to allow dynamic substitution of values](https://golang.org/pkg/text/template/), the syntax is based on double curly braces `{{ expression }}`. The fields supported are `Options`, `Action` and `Meta`, e.g.

`{{ .Options.SourcePath }}` - is the landing zone path  
`{{ .Options.ConfigPath }}` - is the path to the CAF configurations folder
`{{ .Options.StateName }}` - is the name of the state key, plan & state file names and part of DataDir  
`{{ .Options.CafEnvironment }}` - is the source path value  
`{{ .Options.Level }}` - is the value of the level being operated on  
`{{ .Options.Workspace }}` - is the workspace name  
`{{ .Options.DataDir }}` - is path to the data dir (see below)  
`{{ .Options.TargetSubscription }}` - is the Azure subscription ID being deployed into  
`{{ .Options.StateSubscription }}` - is the Azure subscription ID holding state  
`{{ .Options.Subscription.TenantID }}` - Azure tenant ID being used  
`{{ .Options.Identity.ObjectID }}` - is object ID of the signed in identity  
`{{ .Options.Identity.ClientID }}` - is client ID of the signed in identity  
`{{ .Meta.RoverHome }}` - is the path to the rover home directory  

### Grouping commands

RoverGo supports command grouping. This allows you to build a command that is the aggregation of other commands (either builtin or custom commands.)

For example, RoverGo supports plan and apply commands. You create a custom command called deploy to execute both plan and apply.

```yaml
deploy:
  description: "short description"
  commands:
  - plan
  - apply
```

This can be invoked via `rover deploy <options>`

If we wanted to add more commands as part of the deploy group command we can add new commands to the list. 

```yaml
deploy:
  description: "short description"
  commands:
  - plan
  - build
  - lint
  - apply
```

Group commands is a powerful construct that allows commands to be composed into workflows.

## Rover Home Dir

Rover uses `$HOME/.rover` to store data during execution of actions, this is created if it doesn't exist at startup. The default files such as actions.yaml and other configs will be placed there upon creation of the directory

The Rover home directory is also used by Terraform when init is run to hold all of the modules, plugins etc and also the state configuration. This is the **TF_DATA_DIR** and it is set as follows:

`<rover-home>/<workspace>/<level>/<statename>`

Example: If a user dbowie was running Rover using a workspace called "live" and level "level2" and a statename of "web" the **TF_DATA_DIR** would be  

`/home/dbowie/.rover/live/level2/web/`

Within this directory you would expect to see the Terraform modules and providers directories, and also `terraform.tfstate`, and after running a plan the `web.tfplan` file would be here

---

## CAF Concepts

This is not intended to be a complete guide to CAF Landing Zones, which is a complex & nuanced topic, [some of the complete CAF docs can be found here](https://github.com/Azure/caf-terraform-landingzones/tree/master/documentation). However there are a few concepts and terms you will need to understand when running Rover v2

- **Landing zone**
Landing zones (or *CAF Landing Zones*) are simply sets of Azure resources deployed via CAF Terraform modules, they are highly opinionated, governed and designed to support one or more application workloads.

- **Launchpad**
The launchpad is a special type of landing zone which holds Terraform remote state and configuration for all other landing zones. It consists of Azure Storage accounts and Key Vaults. As the launch pad zone holds the state for all other landing zones it has a special lifecycle. It is always deployed with a level identifier of "level0" aka level zero.

- **Level**
Every landing zone is deployed into a level, level 0 is reserved for the launchpad. Levels 1-4 are for normal workload landingzones. **Note**. Even though levels are numeric integers, within CAF they are represented as strings e.g. `level0` and `level1` this is for historical reasons nobody understands.
See the [CAF landing zone docs](https://github.com/Azure/caf-terraform-landingzones/blob/master/documentation/code_architecture/hierarchy.md) for more details on levels

- **Configuration**
When deploying a landing zone you don't write any Terraform directly, but instead provide input configuration and settings in the form of a directory holding `.tfvar` files. This is what Rover refers to as the configuration directory, or simply the configuration (see below)

- **Source Terraform**
The source terraform for the CAF landingzones, in many cases this can be cloned directly from the [https://github.com/Azure/caf-terraform-landingzones](https://github.com/Azure/caf-terraform-landingzones) repo and used unmodified. Rover needs to access this as it contains the Terraform HCL required to deploy a landingzone.

## Integration Tests

Rover has the ability to run integration tests written in Go, using [Terratest.](https://github.com/gruntwork-io/terratest)

When running `rover test` instead of passing a source-dir parameter, a test-source parameter must be passed in. This the path to the tests you would like to execute.

Before running tests, it's expected that resources are deployed using Rover. The pattern for testing is to validate configuration on deployed infrastructure, but the tests themselves do not deploy the infrastructure.

```bash
rover test \
        --config-dir test/testdata/configs/level0/launchpad/ \
        --test-source examples/tests/ \
        --level level0 \
        --environment test \
        --launchpad \
        --statename caf_launchpad
```

Note: --launchpad is only required when testing the level0 launchpad. It is not needed for other landingzones. This is similar behavior to other Rover commands.

The test action also expects a statename parameter. This is the landing zone key the specific landing zone configuration under test. [An example can be found here.](../examples/minimal/configs/level0/launchpad/configuration.tfvars)

Sample tests can be found in the [examples folder of this repo.](../examples/tests)