# Deploy Integration Tests

It's possible to follow the below guidelne to run the deploy integration tests on a fork of `aztfmod/rovergo` repo.

## Guideline

- Fork the `aztfmod/rovergo` repo.

- Create an _Azure Service Principal_ to give permissions to the `aztfmod/rovergo` repo Actions to access the _Azure Subscription_.

```bash
az ad sp create-for-rbac --name "rovergo" --role contributor --sdk-auth
```

> Official documentation can be found at [Create a Service Principal](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli)

- Create a secret in the repo, name it as `AZURE_CREDENTIALS` with the output of _Service Principal_ object from previous step.

- Make a change in the source code, commit and push the changes.

- _CI Builds_ action will kick in automatically, linter, ci builder and integration tester jobs will run sequentially.
