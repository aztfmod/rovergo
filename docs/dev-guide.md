# Rover v2 Internals

Rover v2 makes use of the following:
- [Cobra](https://github.com/spf13/cobra) to provide the core CLI framework and structure of the commands
- [hashicorp/terraform-exec](https://github.com/hashicorp/terraform-exec) for running Terraform
- [Azure SDKs for Go](https://github.com/Azure/azure-sdk-for-go) specifically the ARM APIs, Key Vault and blob storage

## Repo structure and index

```text
\
├── cmd               - Entry points for all the commands and Cobra
├── custom_actions    - Used to define custom actions
├── docs              - You are here
├── pkg               - Go source code for all packages
│   ├── azure         - Wrappers and helpers for calling Azure via the API
│   ├── command       - For calling external commands and exe
│   ├── console       - Console output formatting
│   ├── custom        - Custom actions
│   ├── landingzone   - All code for managing landing zones (more below)
│   ├── symphony      - All code for working with symphony YAML config
│   ├── terraform     - Some terraform helper and handle to tfexec
│   ├── utils         - General stuff
│   └── version       - Version number, common in Go projects to be a package
└── samples           - Reference app configuration
```

## Main concepts, structs types and interfaces

blah

### Action
blah blah

### Options
blah blah blah

## Root cmd

cmd/root.go

## Symphony Code