package azure

import (
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/stretchr/testify/assert"
)

// NOTE. These tests use the Azure CLI currently to get details of the signed in user

func Test_IsOwnerCLI(t *testing.T) {
	// If you're not an owner on the subscription you are using with the az CLI this test will fail
	i := GetSignedInIdentity()
	s := GetSubscription()

	isOwner, err := CheckIsOwner(i.ObjectID, s.ID, s.TenantID)
	assert.Nil(t, err)
	assert.True(t, isOwner)
}

func Test_IsNotOwnerSub(t *testing.T) {
	i := GetSignedInIdentity()
	s := GetSubscription()

	// Random GUID for subscription
	isOwner, err := CheckIsOwner(i.ObjectID, "63f3ec63-61c7-432c-9a10-9513ec3f889e", s.TenantID)
	// This will error with 404
	assert.NotNil(t, err)
	detailedErr := err.(autorest.DetailedError)
	assert.Equal(t, int(404), detailedErr.StatusCode)
	assert.False(t, isOwner)
}

func Test_IsNotOwnerOID(t *testing.T) {
	s := GetSubscription()

	// Random GUID for object id
	isOwner, err := CheckIsOwner("5e441a4a-1da9-4f8e-8022-bd9debec9cc3", s.ID, s.TenantID)
	assert.Nil(t, err)
	assert.False(t, isOwner)
}
