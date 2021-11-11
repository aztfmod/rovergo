# Rover v2 Dev Guide & Internals

Rover v2 makes use of the following:

- [Cobra](https://github.com/spf13/cobra) to provide the core CLI framework and structure of the commands
- [hashicorp/terraform-exec](https://github.com/hashicorp/terraform-exec) for running Terraform
- [Azure SDKs for Go](https://github.com/Azure/azure-sdk-for-go) specifically the ARM APIs, Key Vault and blob storage

## Running & Debugging Locally

Running locally is very simple

```bash
go run main.go {rover CLI switches}
```

You can debug with VS Code, and a [launch.json](./../.vscode/launch.json) configuration is provided, however due to some [limitations](https://github.com/microsoft/vscode/issues/83678) in VS Code, you will need to edit launch.json and place the rover command you wish to debug on into the `args` array section.

## Repo structure and index

```text
\
├── cmd               - Entry points for all the commands and Cobra
├── custom_actions    - Used to define custom actions
├── docs              - You are here
├── pkg               - Go source code for all packages
│   ├── azure         - Wrappers and helpers for calling Azure via the API
│   ├── command       - For calling external commands and capturing output
│   ├── console       - Console output message formatting & logging
│   ├── custom        - Custom actions
│   ├── landingzone   - All code for managing landing zones (more below)
│   ├── symphony      - All code for working with symphony YAML config
│   ├── terraform     - Some terraform helper and handle to tfexec
│   ├── utils         - General stuff ¯\_(ツ)_/¯
│   └── version       - Version number, common in Go projects to be a package
└── examples           - Reference app configuration
```

## Main concepts, structs types and interfaces

The two principal data types & structs in Rover v2 are `Action` and `Options`

### Action

This is an interface defined in `landingzone/actions.go`, as spec it provides no concrete implementation, but is defined as

```go
Execute(o *Options) error
GetName() string
GetDescription() string
```

Use of struct encapsulation and composition is used to provide effectively "sub-classes"

There are many implementations of `Action`,

- `TerraformAction` - found in `landingzones/actions.go`
  - Further extended by `ApplyAction` `PlanAction` etc in the landingzone package, these are concrete implementations which provide the `Execute` function.
- `CustomAction` - found in `custom/action_custom.go`

Actions hold very few fields, aside from a few values specific to their mode of implementation, e.g. CustomActions hold the name of executable they need to call.

The `Execute()` function of the Action is called to actually carry it out.

Actions are used to build the Cobra command structure, see `cmd/root.go` which holds a map of all actions, fixed actions such as TerraformActions and those loaded dynamically at runtime (CustomActions)

#### Terraform Actions

All of the interaction with Terraform for managing CAF landing zones is contained in these six actions. The shared code they all use is held in `pkg/landingzone/landingzone.go`. The actions all follow a general flow of calling `prepareTerraformCAF()` several run `connectToLaunchPad()` and of course enact the relevant terraform command, this is done with tfexec.

The landingzone.go file holds a lot of the shared functions and Rover/CAF specific logic:

- `prepareTerraformCAF()` - This does ***many*** things including getting current details from Azure CLI, getting a tfexec handle (see terraform package), setting many environmental variables (TF_VARS and others), further mutating the Options object preparing it for terraform to run. Also critically it locates the launchpad state storage account for the correct level and environment.

- `connectToLaunchPad()` - Connects to KeyVault and extracts various Terraform variables required by CAF, setting them as TF_VARS

- `runTerraformInit()` - This runs Terraform init in the correct mode, either with remote state and backend configured using the storage account found by `prepareTerraformCAF()` (when NOT running in launchpad mode), or with local state (only used when creating or destroying a launchpad)

### Options

The struct `landingzones.Options` (in pkg/landingzone/options.go) holds all the settings, parameters etc for a Rover operation. It's something of a grab bag but it gathers the large number of parameters Rover requires in one place, combined with the `Action` it defines the execution of Rover when running a command.

`Options` are constructed with two builder functions; `symphony.BuildOptions` and `landingzones.BuildOptions`, one works from parsing the symphony config (see pkg/symphony/parser.go), the other builds the struct from the command line flags (see pkg/landingzone/cobra.go).  
_Note_. In both cases a slice of `Options` is returned, however `landingzones.BuildOptions` will only ever return a single item.

The `SourcePath` and `ConfigPath` fields should never be set directly, instead functions `SetConfigPath()` and `SetSourcePath()` should be used for validation and setting correctly

### Root Cmd

Although the rover executable entrypoint is in main.go, it does almost nothing, the real entry point is `cmd/root.go`. This constructs the main command structure based on Cobra.

Unlike typical Cobra usage where the commands & sub-commands are defined statically, instead a map of Actions is used, and during `init()` this map is used to add sub-commands under the `rootCmd`.

The heart of Rover control flow, is in the root cmd and is small enough to be reproduced here

```go
  action := actionMap[cmd.Name()]

  var optionsList []landingzone.Options

  // Handle symphony mode where config file and level is passed, this will return optionsList with MANY items
  if configFile != "" {
    // Depending on if we're running single or mult-level this will return one or many options
    optionsList = symphony.BuildOptions(cmd)
  }

  // Handle CLI or standalone mode, this will return optionsList with a single item
  if configPath != "" {
    optionsList = landingzone.BuildOptions(cmd)
  }

  for _, options := range optionsList {
    // Now start the action execution...
    // NOTE: If errors occur downstream, depend on logs from there
    err = action.Execute(&options)
    if err != nil {
      os.Exit(1)
    }
  }
```

_Note_. Due to the hybrid/dual-mode of the Rover CLI dependant on the flags provided, very little use of default values for flags has been used, and defaults are handled conditionally in code.

Non-action based commands (e.g. `rover launchpad fetch` are also defined in the cmd package, and using Cobra they also append themselves into the `rootCmd`

### Symphony Code

The file `pkg/symphony/symphony.go` holds the symphony YAML file loader and unmarshaller.

The file `pkg/symphony/parser.go` holds the main parser invoked by `BuildOptions` to parse either all levels or a single level into an slice of `Options`. _Note._ All stacks within a level are always parsed/loaded.

## Dev Tooling

Makefile supports local dev work and CI pipelines, for building, linting and running tests

```text
$ make
build            🔨 Build the rover binary
clean            🧹 Cleanup project
help             💬 This help message :)
lint-fix         🌟 Lint & format, will try to fix errors and modify code
lint             👀 Lint & format, will not fix but sets exit code on error
run              🏃‍ Run locally, with hot reload, it's not very useful
test             🤡 Run unit tests
```

Linting is done with golangci-lint, there is a config file in the project root

### GitHub Actions

Workflows:

- **CI builds** (`.github/workflows/ci-build.yaml`) - Runs on PRs and pushes into the main branch, it runs linting, tests and checks the binary builds ok
- **Release Binaries** (`.github/workflows/release.yaml`) - See below

### Releases

[Goreleaser](https://goreleaser.com/) is used for building and releasing to GitHub, see `.goreleaser.yml`, releases are triggered by pushing git tags with sematic versioning, and run automatically through GitHub Actions, see `.github/workflows/release.yaml`

Binaries are built for Linux, Windows and MacOS

Bash script (install.sh) allows for easy install of binaries direct from GitHub, this was created with [https://github.com/goreleaser/godownloader](https://github.com/goreleaser/godownloader), see the main readme for details on using it.

#### Publish a release

Git tags MUST be in semver for goreleaser, e.g. to publish version `1.2.3`

```bash
git tag 1.2.3
git push origin 1.2.3
```

You can append a pre-release string after the semver digits, e.g. `0.0.1-foobar`

#### Snapshot release

To run or test a local snapshot build, run goreleaser, the results will be place into `./dist/`. This is safe to run anytime. as no GitHub release or tag will be created.

```bash
goreleaser --snapshot --rm-dist
```

## Running Integration Tests

Rovergo has integration tests that you can run locally or via a Github Actions workflow.

### Running integration tests locally

- Clone the `aztfmod/rovergo` repo

```bash
git clone https://github.com/aztfmod/rovergo.git
```

- Clone the Landingzone repo to _rover home directory_

```bash
git clone https://github.com/azure/caf-terraform-landingzones.git ~/.rover/caf-terraform-landingzones
```

- Login to Azure CLI

```bash
az login
```

- Run the test

```bash
 go test ./...  -tags unit -tags integration
```

### Running integration tests via GitHub Actions


- Create an _Azure Service Principal_ to give permissions to the `aztfmod/rovergo` repo Actions to access the _Azure Subscription_.

```bash
az ad sp create-for-rbac --name "rovergo" --role contributor --sdk-auth
```

> Official documentation can be found at [Create a Service Principal](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli)

- Create a secret in the repo, name it `AZURE_CREDENTIALS` with the output of _Service Principal_ object from previous step.
- Create a secret in the repo, name it `ARM_CLIENT_SECRET` with the Service Principal key from the above output.
- Make a change in the source code, commit and push the changes.

- _CI Builds_ action will kick in automatically, linter, ci builder and integration tester jobs will run sequentially.
