# Proposal for usage

```shell
Azure CAF rover is a command line tool in charge of the deployment of the landing zones in your 
Azure environment.
It acts as a toolchain development environment to avoid impacting the local machine but more importantly 
to make sure that all contributors in the GitOps teams are using a consistent set of tools and version.

Usage:
  rover [command (Build-In, Custom or Group)] [Flags]

  Flags:
        --debug     log extra debug information, may contain secrets
    -h, --help      help for rover
    -v, --version   version for rover

Built-In Commands:
  apply       Perform a terraform plan & apply
  plan        Perform a terraform plan
  init        Perform a terraform init and no other action

Custom Commands:
  lint        Perform tflint with list,write parameters
  validate    Perform a terraform validate

Group Commands:
  deploy      deployment workflow for continuous deployment.
              - plan
              - fmt
              - apply
              - validate
              - destroy 

Use "rover [command] --help" for more information about a command.

command file: ~/.rover/commands.yaml
```