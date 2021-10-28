// +build level0

package caf_tests

import (
	"context"
	"fmt"
	"strings"

	"testing"

	"github.com/aztfmod/terratest-helper-caf/state"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

func TestLaunchpadLandingZoneKey(t *testing.T) {
	//arrange
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")

	//act
	landingZoneKey := tfState.GetLandingZoneKey()

	//assert
	assert.Equal(t, "launchpad", landingZoneKey)
}

func TestLaunchpadResourceGroupIsExists(t *testing.T) {
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		name := resourceGroup.GetName()
		exists := azure.ResourceGroupExists(t, name, tfState.SubscriptionID)
		assert.True(t, exists, fmt.Sprintf("Resource group (%s) does not exist", name))
	}
}

func TestLaunchpadResourceGroupIsExistsViaClient(t *testing.T) {
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	client, _ := azure.GetResourceGroupClientE(tfState.SubscriptionID)
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		_, err := client.CheckExistence(context.Background(), rgName)
		assert.NoError(t, err, fmt.Sprintf("Resource group (%s) does not exist", rgName))
	}
}

func TestLaunchpadResourceGroupHasTags(t *testing.T) {
	//arrange
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		level := resourceGroup.GetLevel()

		rg := azure.GetAResourceGroup(t, rgName, tfState.SubscriptionID)

		assert.Equal(t, tfState.Environment, *rg.Tags["environment"], "Environment Tag is not correct")
		assert.Equal(t, "launchpad", *rg.Tags["landingzone"], "LandingZone Tag is not correct")
		assert.Equal(t, level, *rg.Tags["level"], "Level Tag is not correct")
	}
}

func TestLaunchpadResourceGroupHasKeyVault(t *testing.T) {
	//arrange
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		if !strings.Contains(rgName, "security") {
			keyVault, err := tfState.GetKeyVaultByResourceGroup(rgName)
			if err != nil {
				panic(err)
			}

			keyVaultName := keyVault.GetName()

			//act
			kv := azure.GetKeyVault(t, rgName, keyVaultName, tfState.SubscriptionID)

			//assert
			assert.NotNil(t, kv, fmt.Sprintf("KeyVault (%s) does not exists", keyVaultName))
		}
	}
}

func TestLaunchpadKeyVaultHasTags(t *testing.T) {
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		level := resourceGroup.GetLevel()

		if !strings.Contains(rgName, "security") {
			keyVault, err := tfState.GetKeyVaultByResourceGroup(rgName)
			if err != nil {
				panic(err)
			}

			keyVaultName := keyVault.GetName()

			//act
			kv := azure.GetKeyVault(t, rgName, keyVaultName, tfState.SubscriptionID)

			//assert
			assert.NotNil(t, kv, fmt.Sprintf("KeyVault (%s) does not exists", keyVaultName))
			assert.Equal(t, tfState.Environment, *kv.Tags["environment"], "Environment Tag is not correct")
			assert.Equal(t, tfState.Key, *kv.Tags["landingzone"], "LandingZone Tag is not correct")
			assert.Equal(t, level, *kv.Tags["level"], "Level Tag is not correct")
			assert.Equal(t, level, *kv.Tags["tfstate"], "TF State Tag is not correct")
		}
	}
}

func TestLaunchpadResourceGroupHasStorageAccount(t *testing.T) {
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		if !strings.Contains(rgName, "security") {
			storageAccount, err := tfState.GetStorageAccountByResourceGroup(rgName)
			if err != nil {
				panic(err)
			}

			storageAccountName := storageAccount.GetName()

			//act
			storageAccountExists := azure.StorageAccountExists(t, storageAccountName, rgName, tfState.SubscriptionID)

			//assert
			assert.True(t, storageAccountExists, "storage account does not exist")

		}
	}
}

func TestLaunchpadStorageAccountHasTags(t *testing.T) {
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		level := resourceGroup.GetLevel()

		if !strings.Contains(rgName, "security") {
			storageAccount, err := tfState.GetStorageAccountByResourceGroup(rgName)
			if err != nil {
				panic(err)
			}

			storageAccountName := storageAccount.GetName()

			//act
			localStorage, err := azure.GetStorageAccountE(storageAccountName, rgName, tfState.SubscriptionID)

			//assert
			assert.NotNil(t, localStorage, fmt.Sprintf("Storage Account (%s) does not exists", storageAccountName))
			assert.NoError(t, err, "Storage Account couldn't read")
			assert.Equal(t, tfState.Environment, *localStorage.Tags["environment"], "Environment Tag is not correct")
			assert.Equal(t, tfState.Key, *localStorage.Tags["landingzone"], "LandingZone Tag is not correct")
			assert.Equal(t, level, *localStorage.Tags["level"], "Level Tag is not correct")

		}
	}
}

func TestLaunchpadStorageAccountHasTFStateContainer(t *testing.T) {
	t.Parallel()
	tfState := state.NewTerraformState(t, "launchpad")
	resourceGroups := tfState.GetResourceGroups()

	for _, resourceGroup := range resourceGroups {
		rgName := resourceGroup.GetName()
		if !strings.Contains(rgName, "security") {
			storageAccount, err := tfState.GetStorageAccountByResourceGroup(rgName)
			if err != nil {
				panic(err)
			}

			storageAccountName := storageAccount.GetName()
			containerName := "tfstate"

			//act
			exists := azure.StorageBlobContainerExists(t, containerName, storageAccountName, rgName, tfState.SubscriptionID)

			//assert
			assert.True(t, exists, "TF State Container does not exist")

		}
	}
}
