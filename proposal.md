Design and agree the top two (maybe deeper) levels of commands in the new CLI, e.g. `rover landingzone list`

Proposed command structure is
- `landingzone`
  - `fetch` - replaces current --clone operation
  - `run` - Takes an action flag e.g. "plan", "apply", "destry", "validate", "init" etc.
    - `level` specify desired level to deploy to
  - `list` - replaces current --landingzone operation which only support listing
  - `test` - Initiates a terratest execution of the landing zone
- `launchpad` It is felt that this is such a special case it should be it's own command
  - `run` - Same as `landingzone run` but in special launchpad mode
  - `test` - initiates a terratest execution of the launchpad.
- `walkthrough` - Initiate a rover walkthrough to set up an initial project scaffold.
- `workspace` - Wrapper to call terraform workspace commands, is this really required?
- `terraform` - Drop down and run an arbitrary terraform command but under the rover context e.g. logged in user creds etc
- `login` - Will be a stub, I see no value in duplicating az login behaviour
- `logout` - Will be a stub for now, I see no value in duplicating az logout behaviour
- `ci`
  - `{tool}` - name of the specific ci task/tool to run from ci_task definitions
- `cd`

The split between `landingzone` and `launchpad` is a trade off. There will be many shared flags between these commands so on one hand it feels like making launchpad mode a flag e.g. `--launchpad` (as it currently is in rover) but that obscures critical functionally, and giving it space at the top level. Any shared flags and shared implementation (e.g. plan or apply operations) will need to be factored out

Agree on global flags
- `--level-level` Specify the desired log level.
- `--cloud` Name of the specific Azure Cloud.
- `--base-path` Starting path for symphony.
- `--debug` Switch on extra debug info
- `--help`
- `--version`
- `--config` - Loads in rover.yaml config, it remains to be seen if it's needed, it has been supporting the spike/experiment efforts so far

Notes:
* Rovergo should support relative paths and work regardless of where the rover tool is installed.
* Rovergo should not be required to be in /tf/rover.
* Support for Devcontainer and Code Space should be maintained.
* Rovergo should have a complete set of unit test and integration tests.
* Rovergo CI pipeline and PR workflow.

Questions:
* Should we support legacy rover commands or make this a clean break?

See Task #9 for implementing these as stubs so we have a functional CLI with implementation we can plug in as we go along

```bash

# old syntax
rover -lz $CAF_DIR/caf_modules/landingzones/caf_launchpad \
  -launchpad \
  -var-folder $CAF_DIR/config_launchpad/level0/launchpad \
  -parallelism 30 \
  -level level0 \
  -env ${caf_environment} \
  -a apply

# new syntax
rover launchpad run \
      --config $CAF_DIR/configs/level0/launchpad \
      --env demo 

rover launchpad run \
      --path $CAF_DIR/caf_modules/landingzones/caf_launchpad
      --config $CAF_DIR/configs/level0/launchpad \
      --env demo 

rover landingzone run \
      --config $CAF_DIR/configs/level1/foundation \
      --env demo 





rover cd  --config=all/symphony.yml                             # deploy all
rover cd apply --level level0 --config=launchpad/symphony.yml   # deploy level0
rover cd apply --level level1 --config=platform/symphony.yml    # deploy level1
rover cd apply --level level2 --config=platform/symphony.yml    # deploy level1
rover cd apply --level level3 --config=solution/symphony.yml    # deploy level1
rover cd apply --level level4 --config=solution/symphony.yml    # deploy level1

rover cd  --config=all/symphony.yml                               # destroy all
rover cd destroy --level level0 --config=launchpad/symphony.yml   # destroy level0
rover cd destroy --level level1 --config=platform/symphony.yml    # destroy level1
rover cd destroy --level level2 --config=platform/symphony.yml    # destroy level1
rover cd destroy --level level3 --config=solution/symphony.yml    # destroy level1
rover cd destroy --level level4 --config=solution/symphony.yml    # destroy level1

rover cd  --config=all/symphony.yml                            # test all
rover cd test --level level0 --config=launchpad/symphony.yml   # test level0
rover cd test --level level1 --config=platform/symphony.yml    # test level1
rover cd test --level level2 --config=platform/symphony.yml    # test level1
rover cd test --level level3 --config=solution/symphony.yml    # test level1
rover cd test --level level4 --config=solution/symphony.yml    # test level1
```