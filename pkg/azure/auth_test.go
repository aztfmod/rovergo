//go:build integration || ignore || !unit
// +build integration ignore !unit

package azure

import (
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// NOTE. These tests use the Azure CLI currently to get details of the signed in user

func Test_IsOwnerCLI(t *testing.T) {
	// If you're not an owner on the subscription you are using with the az CLI this test will fail
	identity, err := getIdentity()
	assert.Nil(t, err)

	subscription, err := GetSubscription()
	assert.Nil(t, err)

	isOwner, err := CheckIsOwner(identity.ObjectID, subscription.ID)
	assert.Nil(t, err)
	assert.True(t, isOwner)

}

func Test_IsNotOwnerSub(t *testing.T) {
	// arrange
	identity, err := getIdentity()
	assert.Nil(t, err)

	// act
	isOwner, err := CheckIsOwner(identity.ObjectID, utils.GenerateRandomGUID())

	// assert
	assert.NotNil(t, err)
	detailedErr := err.(autorest.DetailedError)
	assert.Equal(t, int(404), detailedErr.StatusCode)
	assert.False(t, isOwner)

}

func Test_IsNotOwnerOID(t *testing.T) {
	s, err := GetSubscription()
	assert.Nil(t, err)

	// Random GUID for object id
	isOwner, err := CheckIsOwner(utils.GenerateRandomGUID(), s.ID)
	assert.Nil(t, err)
	assert.False(t, isOwner)
}

func getIdentity() (*Identity, error) {
	acct, err := GetSubscription()
	if err != nil {
		return nil, err
	}
	switch userType := acct.User.Usertype; userType {
	case "user":
		return GetSignedInIdentity()
	case "servicePrincipal":
		return GetSignedInIdentityServicePrincipal()
	default:
		return nil, nil
	}
}
