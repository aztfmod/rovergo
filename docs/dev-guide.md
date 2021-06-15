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
└── samples           - Reference app configuration
```

## Main concepts, structs types and interfaces

The two principal data types & struts in Rover v2 are `Action` and `Options` 

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

### Options

The struct `landingzones/Options` holds all the settings, parameters etc for a Rover operation. It's something of a grab bag but it gathers everything in one place, combined with the `Action` it defines the execution of Rover when running a command

## Root cmd

cmd/root.go

## Symphony Code

# CI/CD

Makefile supports local dev work and CI pipelines, for building, linting and running tests

```
$ make
build                🔨 Build the rover binary
clean                🧹 Cleanup project
help                 💬 This help message :)
lint-fix             🌟 Lint & format, will try to fix errors and modify code
lint                 👀 Lint & format, will not fix but sets exit code on error
run                  🏃‍ Run locally, with hot reload
test                 🤡 Run tests
```

## Releases

[Goreleaser](https://goreleaser.com/) is used for building and releasing to GitHub, see `.goreleaser.yml`, releases are triggered by pushing git tags with sematic versioning, and run automatically through GitHub Actions, see `.github/workflows/release.yaml`

Binaries are built for Linux, Windows and MacOS

Bash script (install.sh) allows for easy install of binaries direct from GitHub, this was created with https://github.com/goreleaser/godownloader, see the main readme for details on using it.

### Publish a release


### Snapshot release

To run or test a local snapshot build, run goreleaser, the results will be place into `./dist/`. This is safe to run anytime. as no GitHub release or tag will be created.

```bash
goreleaser --snapshot --rm-dist
```