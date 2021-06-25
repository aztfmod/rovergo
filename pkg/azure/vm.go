//
// Rover - Azure Virtual Machine
// * Greg O, June 2021
//

package azure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
)

type Compute struct {
	Name              string
	ResourceGroupName string
}

type Metadata struct {
	Compute Compute
}

func VMInstanceMetadataService() *Metadata {

	client := http.Client{}

	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
	if err != nil {
		return nil
	}

	req.Header.Add("metadata", "true")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
	}

	meta := &Metadata{}
	err = json.Unmarshal([]byte(bodyBytes), meta)
	if err != nil {
		return nil
	}

	return meta
}

// GetVMIdentities will get the MI details of an Azure VM, both system assigned and user assigned
func GetVMIdentities(subID string, resourceGroupName string, vmName string) ([]Identity, error) {
	client := compute.NewVirtualMachinesClient(subID)
	authorizer, err := GetAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = authorizer

	vm, err := client.Get(context.Background(), resourceGroupName, vmName, compute.InstanceViewTypesInstanceView)
	if err != nil {
		return nil, err
	}
	if vm.Identity == nil {
		return nil, nil
	}

	identList := []Identity{}

	if vm.Identity.Type == compute.ResourceIdentityTypeSystemAssigned || vm.Identity.Type == compute.ResourceIdentityTypeSystemAssignedUserAssigned {
		identList = append(identList, Identity{
			DisplayName: "SystemAssigned",
			ObjectID:    *vm.Identity.PrincipalID,
			ClientID:    "UNKNOWN",
			ObjectType:  "servicePrincipal",
		})
	}

	for _, userIdent := range vm.Identity.UserAssignedIdentities {
		identList = append(identList, Identity{
			DisplayName: "UserAssigned",
			ObjectID:    *userIdent.PrincipalID,
			ClientID:    *userIdent.ClientID,
			ObjectType:  "servicePrincipal",
		})
	}

	return identList, nil
}
