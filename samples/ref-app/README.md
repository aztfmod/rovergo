# CAF sample application

### 🔥 NOTE. This has been copied from https://github.com/aztfmod/symphony/tree/master/caf

This CAF application is based on the CAF Terraform Landing Zones Starter Sandpit configuration and used for demonstration purposes.

[caf-terraform-landingzones-starter](https://github.com/Azure/caf-terraform-landingzones-starter/tree/starter/configuration/sandpit)

## Folder structure

One noticable difference between this folder structure and the starter repo is the organization of deployments into three main configurations instead of by level. This structure allows for pipelines to be run by need and permission. There may be multiple config_app_appname repos against this one CAF deployment.

As always, you can adopt your own structure as needed while adhering to the CAF level [principles](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/).

### Landing Zones

- CAF Modules
  - Landing zones imported and unmodified from the main CAF repo
    - [caf-terraform-landingzones](https://github.com/Azure/caf-terraform-landingzones)
    - Imports modules at runtime from the landing zone source defined in the caf_modules/landingzones/caf_*/landingzone.tf
      - [terraform source options](https://www.terraform.io/docs/language/modules/sources.html)
    - Custom code should be in the caf_modules_appname folder (below) to avoid code drift
  - Modules folder for local module reference instead of runtime (as needed)
    - [terraform-azurerm-caf](https://github.com/aztfmod/terraform-azurerm-caf)
    - Use relative path from landing zone source = "../../modules" and no version specification

- CAF Modules APPNAME
  - Landing zones modified for app custom use cases, fixes etc

### Configurations

- Config Launchpad
  - Level 0 only
  - Terraform State management storage containers and keyvaults
  - MSIs for Gitlab runner agents to deploy subsequent levels (moving to an CAF runner add-on eventually)
  - Deployed from the devcontainer rover cli via logged in Owner of target subscription
    - Deploy command available in [local.sh](./local.sh)

- Config Platform
  - Levels 1, 2 and 3
  - Level 3 resources that are shared, not app specific are deployed here (none present in this sample)

- Config App APPNAME
  - Levels 3, 4
  - Configuration for contained application deployment
  - Custom landing zone and modules in matching CAF Modules APPNAME (as needed)
    - This sample uses both caf_modules and caf_modules_argocd to illustrate

## Running with Rover v2

Add `--debug` for more output

You can run all levels at once (by removing the `--level` parameter), but it's saner to run them level by level the first time, 

```bash
rover cd apply --level level0 --symphony-config samples/ref-app/symphony.yml
```
```bash
rover cd apply --level level1 --symphony-config samples/ref-app/symphony.yml
```
```bash
rover cd apply --level level2 --symphony-config samples/ref-app/symphony.yml
```

## Copy to GitLab
Please use the [clone-repos.sh](../scripts/utils/README.md) script to copy all the folders in this path to GitLab as individual repos under a parent Group. This will also add the environmnet variables to the Group to allow proper pipeline and MSI execution.
