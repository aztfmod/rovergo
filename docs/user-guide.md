# Rover v2 User Guide

This is a guide to the Rover v2 command line and operation

## Concepts

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
The source terraform for the CAF landingzones, in many cases this can be cloned directly from the https://github.com/Azure/caf-terraform-landingzones repo and used unmodified. Rover needs to access this as it contains the Terraform HCL required to deploy a landingzone.

# The Rover v2 CLI

```text
Usage:
  rover [command]

Available Commands:
  apply       Perform a terraform plan & apply
  destroy     Perform a terraform destroy
  finder      List all terraform (example custom action)
  fmt         Perform a terraform format
  help        Help about any command
  init        Perform a terraform init and no other action
  landingzone Manage and deploy landing zones
  linter      A linter for terraform
  plan        Perform a terraform plan
  validate    Perform a terraform validate
```

The Rover v2 CLI maps actions (see below) to top level sub commands, with the addition of a special landingzone command

### Action Commands
All commands other than `landingzone` are action commands, they all take the same form and have the same switches, as follows
```
Usage:
  rover <action> [flags]

Flags:
  -v, --config-dir string    Configuration directory, you must supply this or config-file
  -c, --config-file string   Configuration file, you must supply this or config-dir
  -d, --dry-run              Execute a dry run where no actions will be executed
  -e, --environment string   Name of CAF environment
  -h, --help                 help for init
      --launchpad            Run in launchpad mode, i.e. level0
  -l, --level string         CAF landingzone level name, default is all levels
  -s, --source string        Path to source of landingzone
      --state-sub string     Azure subscription ID where state is held
  -n, --statename string     Name for state and plan files
      --target-sub string    Azure subscription ID to operate on
  -w, --workspace string     Name of workspace

Global Flags:
      --debug                Log extra debug information, may contain secrets
```

### Landingzone Management Commands
```
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
 - **Custom actions**, these extend the Rover v2 command set, and allow you to run an external command against the given CAF config, level, source etc.

# Running Rover v2

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
rover apply --config-dir ./caf-config/level1/myapp --source ./landingzones --level level1
```

- Running apply action for a landing zone with config held in `level1/myapp` this time we deploy with the CAF environment set to "prod"
```bash
rover apply --config-dir ./caf-config/level1/myapp --source ./landingzones --level level1 --environment prod
```

## Running in "Config file mode" (multi level)
This mode is intended for use in CI/CD pipelines when multiple levels are being managed at once, and each of those levels contains multiple "stacks". All settings are held within a YAML configuration file (which is [part of project symphony](https://github.com/aztfmod/symphony)). When running with Rover either a single level is specified or all levels are run

For an example and reference symphony config file see [samples/ref-app-symphony.yaml](./../samples/ref-app-symphony.yaml)

**💬 NOTE. This mode is engaged by the use of the `--config-file` switch**

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

### Shared switches
- `--level` Set which level is being operated on
- `--dry-run` Set to perform a dry run and output details of the operation without executing it.

### Ad-hoc mode switches

- `--source` The source landingzone repo location
- `--launchpad` Run deployment in launchpad mode, *only* run this with a valid launchpad config and with level set to "level0"
- `--environment` Set the CAF environment, **defaults to "sandpit"**
- `--statename` Set the state name used for naming tfstate files, **defaults to "caf_solution" or "caf_launchpad"** depending if --launchpad is set
- `--target-sub` Subscription ID to deploy resources into, defaults to current set on Azure CLI
- `--state-sub` Subscription ID where state (i.e. the launchpad) is held, defaults to current set on Azure CLI
- `--workspace` Workspace is used to name the containers used for state, **defaults to "tfstate"**

# Custom Actions

## Rover home
