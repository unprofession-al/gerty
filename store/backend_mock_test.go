package store

import (
	"os"
	"testing"

	"github.com/unprofession-al/gerty/entities"
)

var store Store

func TestMain(m *testing.M) {
	roleMockData := map[int64]*RoleMock{
		1: {Name: "ra", Vars: "[{\"name\":\"Bucket 1\",\"prio\":1,\"vars\":{\"key_1\": \"val_a\"}}]", Children: []int64{2, 3}},
		2: {Name: "rb"},
		3: {Name: "rc"},
	}
	roleBackendMock := NewRoleBackendMock(roleMockData)

	nodeMockData := []*NodeMock{
		{Name: "Node 1", Vars: "[{\"name\":\"Bucket X\",\"prio\":1,\"vars\":{\"key_1\": \"val_node\"}}]", Roles: []int64{2, 3}},
	}
	nodeBackendMock := NewNodeBackendMock(nodeMockData, &roleBackendMock)

	store = Store{
		Roles: &roleBackendMock,
		Nodes: &nodeBackendMock,
	}

	os.Exit(m.Run())
}

func TestInitRoleMock(t *testing.T) {
	rc, _ := store.Roles.Get("rc")
	if rc.Depth() != 1 {
		t.Error("RoleBackendMock seems not to be initialized properly")
	}
}

func TestInitNodeMock(t *testing.T) {
	node, _ := store.Nodes.Get("Node 1")
	found := []string{}
	requiredRoles := []string{"rb", "rc"}
	for _, requiredRole := range requiredRoles {
		for _, role := range node.Roles {
			if requiredRole == role.Name {
				found = append(found, role.Name)
			}
		}
	}

	if len(requiredRoles) != len(found) {
		t.Errorf("found roles %v, should have %v", found, requiredRoles)
	}
}

func TestGetNotExistingNode(t *testing.T) {
	node, err := store.Nodes.Get("Node does not exist")
	if err == nil {
		t.Errorf("node does exist but should not: %v", node)
	}
}

func TestSaveToNodeMock(t *testing.T) {
	name := "Test Node"
	role, _ := store.Roles.Get("rc")
	roles := entities.Roles{}
	roles = append(roles, &role)

	test := entities.Node{
		Name:  name,
		Roles: roles,
		Vars:  entities.VarCollection{},
	}
	store.Nodes.Save(test)

	node, _ := store.Nodes.Get(name)
	for _, nodeRole := range node.Roles {
		if role.Name != nodeRole.Name {
			t.Errorf("node seems not to be saved properly")
		}
	}
}
