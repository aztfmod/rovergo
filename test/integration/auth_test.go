package test

import (
	"testing"
	"time"

	"github.com/aztfmod/rover/cmd"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/landingzone"
	rovertesting "github.com/aztfmod/rover/pkg/testing"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_VM_No_ID(t *testing.T) {
	//
	// this one doesn't test anything for rover. It's here as a trivial template.
	// expectation going in is the VM has no identities
	//
	_, err := rovertesting.AzLogin(t, "-i")

	assert.Error(t, err)

}

func TestIntegration_VM_SystemAssigned_No_Role(t *testing.T) {
	//
	// this one doesn't test anything for rover. It's here as a trivial template.
	//
	defer clearup(t)

	// use bootstrap identity to get the ball rolling
	_, err := rovertesting.AzLoginBootstrap(t)
	if err != nil {
		t.Fatal(err)
	}

	// add the system assigned MI with no role assignments
	_, err = rovertesting.AzVMIdentityAssign(t, "[system]", "")
	assert.NoError(t, err)

	// log in as the system assigned MI
	_, err = rovertesting.AzLogin(t, "-i")

	// should error out because no role assignment & no --allow-no-subscriptions for tenant level access
	assert.Error(t, err)

}

func TestIntegration_VM_SystemAssigned_SubOwner_Role(t *testing.T) {

	defer clearup(t)

	// use bootstrap identity to get the ball rolling
	_, err := rovertesting.AzLoginBootstrap(t)
	if err != nil {
		t.Fatal(err)
	}

	// add the system assigned MI with subscription owner role assignment
	_, err = rovertesting.AzVMIdentityAssign(t, "[system]", "Owner")
	if err != nil {
		t.Fatal(err)
	}

	// logout
	err = rovertesting.AzLogout(t)
	if err != nil {
		t.Fatal(err)
	}

	loginSuccessful := false
	for i := 0; i < 20; i++ {

		// log in as the system assigned MI
		_, err = rovertesting.AzLogin(t, "-i")
		if err == nil {
			loginSuccessful = true
			break
		}

		console.Warning("Waiting 15 seconds for next attempt")
		time.Sleep(time.Second * 15)

	}
	if loginSuccessful == false {
		t.Fatal("Failed to login as system assigned ID within 5 minutes")
	}

	// get the object id of the system assigned MI
	vmIdentityDetails, err := rovertesting.AzVMIdentityShow(t)
	if err != nil {
		t.Fatal(err)
	}

	// set up a terraform fmt command for the actual test
	testCmd := &cobra.Command{
		Use: "fmt",
	}
	testCmd.Flags().Bool("dry-run", true, "")
	testCmd.Flags().String("config-dir", "../testdata/configs/level0/launchpad", "")
	testCmd.Flags().String("source", "../testdata/caf-terraform-landingzones", "")
	testCmd.Flags().String("level", "level0", "")
	testCmd.Flags().Bool("launchpad", true, "")

	optionsList := landingzone.BuildOptions(testCmd)

	action := cmd.ActionMap[testCmd.Name()]
	_ = action.Execute(&optionsList[0])

	assert.Equal(t, "servicePrincipal", optionsList[0].Identity.ObjectType)
	assert.Equal(t, "SystemAssigned", optionsList[0].Identity.DisplayName)
	assert.Equal(t, vmIdentityDetails.PrincipalID, optionsList[0].Identity.ObjectID)
}

func clearup(t *testing.T) {
	err := rovertesting.AzLogout(t)
	if err != nil {
		t.Fatal(err)
	}

	_, err = rovertesting.AzLoginBootstrap(t)
	if err != nil {
		t.Fatal(err)
	}

	vmIdentityDetails, err := rovertesting.AzVMIdentityShow(t)
	if err != nil {
		t.Fatal(err)
	}

	ownerRoleAssignmentID, err := getOwnerRoleAssignmentID(t, vmIdentityDetails.PrincipalID)
	if err != nil {
		t.Fatal(err)
	}

	if ownerRoleAssignmentID != "" {
		err = rovertesting.AzRoleAssignmentDelete(t, ownerRoleAssignmentID)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = rovertesting.AzVMIdentityRemove(t, "[system]")
	if err != nil {
		t.Fatal(err)
	}

	err = rovertesting.AzLogout(t)
	if err != nil {
		t.Fatal(err)
	}
}

func getOwnerRoleAssignmentID(t *testing.T, principalID string) (string, error) {

	roleAssignments, err := rovertesting.AzRoleAssignmentList(t)
	if err != nil {
		t.Fatal(err)
	}

	for _, roleAssignment := range roleAssignments {

		if roleAssignment.PrincipalID == principalID {
			return roleAssignment.ID, nil
		}
	}

	return "", nil
}
