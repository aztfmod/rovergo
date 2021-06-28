package testing

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/joho/godotenv"
)

const OwnerBuiltInRole string = "8e3af657-a8ff-443c-a75c-2fe8c4bcb635"

type VMIdentityShow struct {
	PrincipalID            string                     `json:"principalId,omitempty"`
	TenantID               string                     `json:"tenantId,omitempty"`
	ResourceType           string                     `json:"type,omitempty"`
	UserAssignedIdentities map[string]json.RawMessage `json:"userAssignedIdentities,omitempty"`
}

type IdentityAssignment struct {
	SystemAssignedIdentity string                     `json:"systemAssignedIdentity,omitempty"`
	UserAssignedIdentities map[string]json.RawMessage `json:"userAssignedIdentities,omitempty"`
}

type RoleAssignment struct {
	CanDelegate      string `json:"canDelegate,omitempty"`
	Condition        string `json:"condition,omitempty"`
	ConditionVersion string `json:"conditionVersion,omitempty"`
	Description      string `json:"description,omitempty"`
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	PrincipalID      string `json:"principalId,omitempty"`
	PrincipalType    string `json:"principalType,omitempty"`
	RoleDefinitionID string `json:"roleDefinitionId,omitempty"`
	Scope            string `json:"scope,omitempty"`
	ResourceType     string `json:"type,omitempty"`
}

// UserAssignedIdentity is the output of az identity create, i.e. the details of a newly created user assigned identity
type UserAssignedIdentity struct {
	ClientID        string                     `json:"clientID,omitempty"`
	ClientSecretURL string                     `json:"clientSecretUrl,omitempty"`
	ID              string                     `json:"id,omitempty"`
	Location        string                     `json:"location,omitempty"`
	Name            string                     `json:"name,omitempty"`
	PrincipalID     string                     `json:"principalId,omitempty"`
	ResourceGroup   string                     `json:"resourceGroup,omitempty"`
	Tags            map[string]json.RawMessage `json:"tags,omitempty"`
	TenantID        string                     `json:"tenantId,omitempty"`
	ResourceType    string                     `json:"type,omitempty"`
}

type TestConfig struct {
	SubscriptionID      string `yaml:"subscriptionID,omitempty"`
	VMResourceGroupName string `yaml:"vmResourceGroupName,omitempty"`
	VMName              string `yaml:"vmName,omitempty"`
	Location            string `yaml:"location,omitempty"`
	SPNUsername         string `yaml:"spnUsername,omitempty"`
	SPNPassword         string `yaml:"spnPassword,omitempty"`
	TenantID            string `yaml:"tenantID,omitempty"`
}

type SPNDetails struct {
	DisplayName           string
	ObjectID              string
	ServicePrincipalNames []string
}

// Config is exported for some reason?
var Config TestConfig

func NewTestConfiguration() (*TestConfig, error) {

	tc := new(TestConfig)

	err := godotenv.Load("testConfig.env")
	if err != nil {
		console.Error("could not load testConfig.env")
		return nil, err
	}

	tc.Location = os.Getenv("ROVER_LOCATION")
	tc.SPNPassword = os.Getenv("ROVER_SPNPASSWORD")
	tc.SPNUsername = os.Getenv("ROVER_SPNUSERNAME")
	tc.SubscriptionID = os.Getenv("ROVER_SUBSCRIPTIONID")
	tc.TenantID = os.Getenv("ROVER_TENANTID")
	tc.VMName = os.Getenv("ROVER_VMNAME")
	tc.VMResourceGroupName = os.Getenv("ROVER_VMRESOURCEGROUPNAME")

	return tc, nil
}

func init() {

	tmpConfig, err := NewTestConfiguration()
	if err != nil {
		console.Error("could not load testConfig.env")
	}
	Config = *tmpConfig
}

func AzVMIdentityAssign(t *testing.T, identity string, role string) (*IdentityAssignment, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "vm", "identity", "assign"}
	args = append(args, []string{"--identities", identity}...)

	// if the assignment target is a user assigned identity, it already has a role assignment
	// if the assignment target is a system assigned identity, it may receive owner or it may receive no perms
	if identity == "[system]" && strings.EqualFold(role, "owner") {
		args = append(args, []string{"--role", "Owner"}...)
		args = append(args, []string{"--scope", fmt.Sprintf("/subscriptions/%s", Config.SubscriptionID)}...)
	}
	args = append(args, []string{"--name", Config.VMName}...)
	args = append(args, []string{"--resource-group", Config.VMResourceGroupName}...)

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	identityAssignment := &IdentityAssignment{}
	err = json.Unmarshal([]byte(cmdRes), identityAssignment)
	if err != nil {
		return nil, err
	}

	return identityAssignment, nil
}

func AzRoleAssignmentCreate(t *testing.T, assigneeObjectID string) (*RoleAssignment, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "role", "assignment", "create"}
	args = append(args, []string{"--role", OwnerBuiltInRole}...)
	args = append(args, []string{"--assignee-object-id", assigneeObjectID}...)
	args = append(args, []string{"--assignee-principal-type", "ServicePrincipal"}...)
	args = append(args, []string{"--scope", fmt.Sprintf("/subscriptions/%s", Config.SubscriptionID)}...)

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	roleAssignment := &RoleAssignment{}
	err = json.Unmarshal([]byte(cmdRes), roleAssignment)
	if err != nil {
		return nil, err
	}

	return roleAssignment, nil
}

func AzRoleAssignmentDelete(t *testing.T, roleAssignmentID string) error {
	err := command.CheckCommand("az")
	if err != nil {
		return err
	}

	args := []string{"az", "role", "assignment", "delete"}
	args = append(args, []string{"--ids", roleAssignmentID}...)

	_, err = command.QuickRun(args...)
	if err != nil {
		return err
	}

	return nil
}

func AzRoleAssignmentList(t *testing.T) ([]RoleAssignment, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "role", "assignment", "list", "--role", "owner"}

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	roleAssignments := &[]RoleAssignment{}
	err = json.Unmarshal([]byte(cmdRes), roleAssignments)
	if err != nil {
		return nil, err
	}

	return *roleAssignments, nil
}

func AzIdentityDelete(t *testing.T, identityID string) error {
	err := command.CheckCommand("az")
	if err != nil {
		return err
	}

	args := []string{"az", "identity", "delete"}
	args = append(args, []string{"--ids", identityID}...)

	_, err = command.QuickRun(args...)
	if err != nil {
		return err
	}

	return nil
}

func AzIdentityCreate(t *testing.T, identityName string) (*UserAssignedIdentity, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "identity", "create"}
	args = append(args, []string{"--name", identityName}...)
	args = append(args, []string{"--resource-group", Config.VMResourceGroupName}...)
	args = append(args, []string{"--location", Config.Location}...)

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	userAssignedIdentity := &UserAssignedIdentity{}
	err = json.Unmarshal([]byte(cmdRes), userAssignedIdentity)
	if err != nil {
		return nil, err
	}

	console.Debugf("New user assigned identity %s created.", userAssignedIdentity.Name)

	return userAssignedIdentity, nil
}

func AzVMIdentityRemove(t *testing.T, identityName string) error {
	err := command.CheckCommand("az")
	if err != nil {
		return err
	}

	args := []string{"az", "vm", "identity", "remove"}
	args = append(args, []string{"--identities", identityName}...)
	args = append(args, []string{"--resource-group", Config.VMResourceGroupName}...)
	args = append(args, []string{"--name", Config.VMName}...)

	_, err = command.QuickRun(args...)
	if err != nil {
		return err
	}

	console.Debugf("Identity %s removed.", identityName)

	return nil
}

func AzVMIdentityShow(t *testing.T) (*VMIdentityShow, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "vm", "identity", "show"}
	args = append(args, []string{"--resource-group", Config.VMResourceGroupName}...)
	args = append(args, []string{"--name", Config.VMName}...)

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	vmIdentityShow := &VMIdentityShow{}
	if cmdRes != "" {
		err = json.Unmarshal([]byte(cmdRes), vmIdentityShow)
		if err != nil {
			return nil, err
		}
	}

	return vmIdentityShow, nil
}

func AzLogin(t *testing.T, parms ...string) (*azure.Subscription, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "login"}
	args = append(args, parms...)

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	subs := &[]azure.Subscription{}
	err = json.Unmarshal([]byte(cmdRes), subs)
	if err != nil {
		return nil, err
	}

	sub := &(*subs)[0]

	console.Debugf("Azure subscription is: %s (%s)\n", sub.Name, sub.ID)
	console.Debugf("Logged in security user: %s (%s). Identified by: %s  \n", sub.User.Name, sub.User.Usertype, sub.User.AssignedIdentityInfo)
	return sub, nil
}

func AzLogout(t *testing.T) error {
	err := command.CheckCommand("az")
	if err != nil {
		return err
	}

	args := []string{"az", "logout"}

	_, err = command.QuickRun(args...)
	if err != nil {
		return err
	}

	return nil
}

func AzLoginBootstrap(t *testing.T) (*azure.Subscription, error) {
	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "login", "--service-principal"}
	args = append(args, []string{"--username", Config.SPNUsername}...)
	args = append(args, []string{"--password", Config.SPNPassword}...)
	args = append(args, []string{"--tenant", Config.TenantID}...)

	_, err = command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	return &azure.Subscription{}, nil
}

func AzADSPShow(t *testing.T, appid string) (*SPNDetails, error) {

	err := command.CheckCommand("az")
	if err != nil {
		return nil, err
	}

	args := []string{"az", "ad", "sp", "show"}
	args = append(args, []string{"--id", appid}...)

	cmdRes, err := command.QuickRun(args...)
	if err != nil {
		return nil, err
	}

	spn := &SPNDetails{}
	err = json.Unmarshal([]byte(cmdRes), spn)
	if err != nil {
		return nil, err
	}

	return spn, nil
}
