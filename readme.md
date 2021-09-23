# 🐶 Rover v2

The Rover v2 project undertakes re-writing the [Rover command line tool](https://github.com/aztfmod/rover) in Go and redesigning the running and operation of the tool.

The high level goals & objectives of this project:

- Move away from a bash script based tool to one with solid engineering.
- Provide a maintainable codebase going forward, supported by CI/CD and tests etc.
- Improve the user experience, and simplify the getting started with CAF landingzones.
- Remove dependency on hardcoded paths and other limitations that Rover v1 had.
- Provide standalone binaries (Linux, Windows, MacOS) as well as devcontainer support

Note. This version is not backwards compatible with the previous rover v1 tool

## 🥇 Intro To Rover

Rover is a command line tool to assist with the deployment and management of Azure CAF Landing zones. It provides a way to run Terraform and other tools in a structured way and simplify many of the operations. It was designed to be run by end users locally, but also inside of a CI/CD pipeline

Rover v2 provides a way to manage entire environments consisting of any number of CAF Landing zones across multiple CAF levels, this is done with YAML definitions describing your configuration (aka [project Symphony](https://github.com/aztfmod/symphony)) this is how it used

## 🚦 Project Status

![last commit](https://img.shields.io/github/last-commit/aztfmod/rovergo)
![commit activity](https://img.shields.io/github/commit-activity/w/aztfmod/rovergo)
![release](https://img.shields.io/github/release/aztfmod/rovergo)
![checks status](https://img.shields.io/github/checks-status/aztfmod/rovergo/main)
![ci build status](https://img.shields.io/github/workflow/status/aztfmod/rovergo/CI%20builds?label=ci-build)
![workflow status](https://img.shields.io/github/workflow/status/aztfmod/rovergo/Release%20Binaries?label=release)
![license](https://img.shields.io/github/license/aztfmod/rovergo)

### Board: [aztfmod/projects](https://github.com/orgs/aztfmod/projects/28?card_filter_query=label%3Arover-go)

### Estimated stability level: ⬛⬛⬛⬛⬛⬛⬛⬛⬜⬜ 85%

---

## ✨ Getting Started

### 📦 Installation

The easiest way to install Rover v2 is directly from GitHub using the install script.

To install into the current directory `./bin/` directory just run:

```bash
curl https://raw.githubusercontent.com/aztfmod/rovergo/main/install.sh | bash
```

Alternatively specify the install location, e.g. a directory on your system path:

```bash
export installPath=$HOME/.local/bin
curl https://raw.githubusercontent.com/aztfmod/rovergo/main/install.sh | bash -s -- -b $installPath
```

If you wish to fetch the binary yourself or install older versions, please refer to [the GitHub releases page](https://github.com/aztfmod/rovergo/releases)

### 🏃‍♂️ Running

This is a very basic introduction into running Rover v2

Firstly Rover v2 relies on some external tools and dependencies:

- Azure CLI for authentication and sign-in to Azure - [Install here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- Terraform v0.15+ for deployment of landing zones - [Install here](https://www.terraform.io/downloads.html)

If you want to get up and running and try Rover out, please [check out the minimal sample config](./examples/minimal/readme.md) which has a very simple sample CAF configuration to use, consisting of a level 0 launchpad and a level 1 with some simple resource groups for testing.

A more complex sample is found in the [reference sample app](./examples/reference/readme.md)

Otherwise:

- You will also require a [set of CAF Landingzones](https://github.com/Azure/caf-terraform-landingzones) on your system
- In addtion you will need a set of of landing zone configurations, this is a complex topic beyond the scope of this readme, see the [full user guide](docs/user-guide.md) for some details.

Run rover with `rover --help` to get information about the commands available, e.g.

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
  lint        Run tflint
  plan        Perform a terraform plan
  validate    Perform a terraform validate
```

#### [See the user guide for more details on how to use rover](./docs/user-guide.md)

### 🔌 Extending Rover

Rover v2 is extensible with custom actions which extend it beyond running Terraform, see the [user guide for details](docs/user-guide.md)

### 👩‍💻 Developing and Contributing

#### [See the developer guide](./docs/dev-guide.md)
