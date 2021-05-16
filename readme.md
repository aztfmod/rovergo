# Rover Go

A proof of concept re-writing Rover in Go

Uses [Cobra](https://github.com/spf13/cobra) to provide the framework for a robust and familiar CLI tool and [Viper](https://github.com/spf13/viper) for configuration

# Current status

### ☢ This should be considered spike / poc / investigation grade code

The rover "clone" command has been implemented as a test of the effort.
External command integration & execution is done via `pkg/command/command.go`

The "landingzone" command has been split into "run" and "list" sub-commands
The "run" sub-command currently can perform init, plan and apply and also connect to remote state. It uses [hashicorp/terraform-exec](https://github.com/hashicorp/terraform-exec) to interface with the terraform CLI

The `.rover.yaml` config file is used for all config and also holds remote state (backend) details for Terraform, as this includes storage access keys it's kept out of git. See the `.rover.sample.yaml`

This file needs to be located where `rover` is executing or in your home directory

# Example outputs from the CLI

```
$ ./rover

Azure CAF rover is a command line tool in charge of the deployment of the landing zones in your
Azure environment.
It acts as a toolchain development environment to avoid impacting the local machine but more importantly
to make sure that all contributors in the GitOps teams are using a consistent set of tools and version.

Usage:
  rover [command]

Available Commands:
  clone       Fetch supporting artifacts such as landingzones from GitHub
  completion  Generate completion script
  help        Help about any command
  landingzone Manage landing zones

Flags:
      --config string   config file (default is $HOME/.rover.yaml)
  -h, --help            help for rover
  -t, --toggle          Help message for toggle
  -v, --version         version for rover

Use "rover [command] --help" for more information about a command.
```

Landingzone Run command

```
Run an action for landingzones

Usage:
  rover landingzone run [flags]

Flags:
  -a, --action string      Action to run, one of [plan | apply | destroy] (default "plan")
  -e, --env string         Name of environment (default "caf")
  -h, --help               help for run
  -l, --level int          Level (default 1)
  -s, --source string      Source of landingzone (required)
  -n, --statename string   State and plan name (default "mystate")
  -v, --vars string        Where configuration vars are located (default ".")

Global Flags:
      --config string   config file (default is $HOME/.rover.yaml)
```

Clone command

```
$ ./rover clone --help
Pull down repos from GitHub and extracts them in well defined way.
Git is not required

Usage:
  rover clone [flags]

Flags:
  -b, --branch string   Which branch to clone (default "master")
  -d, --dest string     Where to place output (default "./landingzones")
  -f, --folder string   Extract a sub-folder from the repo
  -h, --help            help for clone
  -r, --repo string     Which repo on GitHub to clone (default "azure/caf-terraform-landingzones")
  -s, --strip int       Levels to strip from repo hierarchy, best left as 1 (default 1)

Global Flags:
      --config string   config file (default is $HOME/.rover.yaml)

```

![](https://user-images.githubusercontent.com/14982936/118290956-edc05680-b4ce-11eb-9b08-409f3bc8679c.png)
