# Getting Started Config for Rover v2

This is an minimal working configuration set for testing Rover v2 with.  
It is intended for first time users, sanity checking and for Rover developers to test with (integration/smoke tests)

## Usage

From the rovergo project repo root, cloned on your system. First fetch the landingzones source.

```bash
rover landingzone fetch
```

Deploy the launchpad but check everything is OK first with a dry run

```bash
rover apply -c examples/minimal/symphony.yaml -l level0 --dry-run
```

Now deploy the launchpad for real (no dry run)

```bash
rover apply -c examples/minimal/symphony.yaml -l level0
```

Deploy level 1

```
rover apply -c examples/minimal/symphony.yaml -l level1
```

### Level 0 - configs/level0/launchpad

This is a copy of the caf-starter demo level0 config for a launchpad
It expects that workspace is set to `tfstate`

### Level 1 - configs/level1

This config contains just two "stacks" (web and test) each with a single Azure resource group.
It is intended to provide the absolute bare minimum deployable set of resources to test Rover v2 with
