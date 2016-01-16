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
		4: {Name: "rd"},
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

func TestSaveNewRole(t *testing.T) {
	rx, _ := store.Roles.Get("rc")
	rx.Name = "rx"
	err := store.Roles.Save(rx)

	rtest, err := store.Roles.Get("rx")
	if err != nil {
		t.Error("Role could net be added prorpely")
	}

	if rx.Parent != rtest.Parent {
		t.Error("Role could net be added prorpely")
	}
}

func TestDeleteExistingRole(t *testing.T) {
	name := "rd"
	role, _ := store.Roles.Get(name)
	err := store.Roles.Delete(role)
	if err != nil {
		t.Errorf("role seems not to be deleted properly")
	}

	_, err = store.Roles.Get(name)
	if err == nil {
		t.Errorf("role seems not to be deleted properly")
	}
}

func TestDeleteInexistingRole(t *testing.T) {
	test := entities.Role{
		Name: "Inexisting Role",
	}
	err := store.Roles.Delete(test)
	if err == nil {
		t.Errorf("Role seems to be deleted but should not")
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

func TestGetInexistingNode(t *testing.T) {
	node, err := store.Nodes.Get("Node does not exist")
	if err == nil {
		t.Errorf("node does exist but should not: %v", node)
	}
}

func TestSaveNewNode(t *testing.T) {
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

func TestSaveExistingNode(t *testing.T) {
	name := "Node 1"
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

func TestDeleteExistingNode(t *testing.T) {
	name := "Node 1"
	node, _ := store.Nodes.Get(name)
	err := store.Nodes.Delete(node)
	if err != nil {
		t.Errorf("node seems not to be deleted properly")
	}

	node, err = store.Nodes.Get(name)
	if err == nil {
		t.Errorf("node seems not to be deleted properly")
	}
}

func TestDeleteInexistingNode(t *testing.T) {
	test := entities.Node{
		Name: "Inexisting Node",
	}
	err := store.Nodes.Delete(test)
	if err == nil {
		t.Errorf("node seems to be deleted but should not")
	}
}
