# Launchpad Notes - Current Implementation

```bash
rover -lz /tf/caf/landingzones/caf_launchpad \
  -launchpad \
  -var-folder /tf/caf/configuration/${environment}/level0/launchpad \
  -parallelism 30 \
  -level level0 \
  -env ${caf_environment} \
  -a [plan|apply|destroy]
```

Initial set up
  - `landingzone_name` = path to TF modules directory, i.e. landingzones
  - `TF_VAR_tf_name` = last part of path OR user specify with -tfstate flag
  - `caf_command`="launchpad" this enables LP / Launchpad mode
  - `tf_action`= what was passed to -a i.e. [ plan | apply | destroy ]
  - `TF_VAR_workspace`="tfstate"
  
## 🔥 MAJOR TODO! verify login scenario and behavior IS AZ LOGIN ALWAYS USED???

- Call fn - *verify_azure_session()*
  - runs `az account show`
- Call fn - *process_target_subscription()*
  - if -target_subscription flag passed then `az account set` to that subid
  - `target_subscription_name` & `target_subscription_id` obtained from `az account` command
  - `ARM_SUBSCRIPTION_ID`=`target_subscription_id`
  - `TF_VAR_tfstate_subscription_id` same as `ARM_SUBSCRIPTION_ID` *unless overriden* with -tfstate_subscription_id flag
  - In LP mode if `TF_VAR_tfstate_subscription_id` != `ARM_SUBSCRIPTION_ID` then EXIT!

Gets terraform_version  for reason?

Call fn *verify_parameters()*
  - If in LP mode && `landingzone_name` UNSET -> Display error and exit
  - Ensures `TF_VAR_tf_name` and `TF_VAR_tf_plan` are set AGAIN to last segment of `landingzone_name` with .tfstate and .tfplan appended.

Call fn *deploy()* -> pass in TF_VAR_workspace
  Call fn *get_storage_id()*
    - Try to find a storage acct in TF_VAR_tfstate_subscription_id WHERE tags match level and environment and set `id` to it
  remove local state and plan
  check if destroy and no storage acct -> error
  Call fn *initialize_state()*
    - Check we are owner on subscription otherwise error and exit
    - Remove sudo rm -f -- ${landingzone_name}/backend.azurerm.tf I DON@T KNOW WHY!
    - remove any local state files
    - set TF_VAR_tf_name and TF_VAR_tf_plan AGAIN!
    - Call `terraform init -upgrade=true` in the landingzone_name folder
    - then calls plan/apply/destroy as per `tf_action`