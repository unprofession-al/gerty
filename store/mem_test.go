package store

import (
	"testing"

	"github.com/unprofession-al/gerty/entities"
)

func TestRoleStoreMemSave(t *testing.T) {
	ri := RoleStoreMem{roles: make(map[string]entities.Role)}
	newRole := entities.Role{Name: "TestRole"}

	ri.Save(newRole)

	role, _ := ri.Get(newRole.Name)

	if role.Name != newRole.Name {
		t.Errorf("Name in not consistent: `%s` != '%s'", role.Name, newRole.Name)
	}
}

func TestRoleStoreMemDelete(t *testing.T) {
	ri := RoleStoreMem{roles: make(map[string]entities.Role)}
	newRole := entities.Role{Name: "TestRole"}

	ri.Save(newRole)
	ri.Delete(newRole)

	_, err := ri.Get(newRole.Name)

	if err == nil {
		t.Errorf("Role '%s' was not properly deleted", newRole.Name)
	}
}

func TestRoleStoreMemList(t *testing.T) {
	ri := RoleStoreMem{roles: make(map[string]entities.Role)}
	role1 := entities.Role{Name: "TestRole1"}
	role2 := entities.Role{Name: "TestRole2"}

	ri.Save(role1)
	ri.Save(role1)
	ri.Save(role2)

	list := ri.List()

	if len(list) != 2 {
		t.Errorf("Roles where not listed properly")
	}
}

func TestNodeStoreMemSave(t *testing.T) {
	ni := NodeStoreMem{nodes: make(map[string]entities.Node)}
	newNode := entities.Node{Name: "TestNode"}

	ni.Save(newNode)

	node, _ := ni.Get(newNode.Name)

	if node.Name != newNode.Name {
		t.Errorf("Name in not consistent: `%s` != '%s'", node.Name, newNode.Name)
	}
}

func TestNodeStoreMemDelete(t *testing.T) {
	ni := NodeStoreMem{nodes: make(map[string]entities.Node)}
	node := entities.Node{Name: "TestNode"}

	ni.Save(node)
	ni.Delete(node)

	_, err := ni.Get(node.Name)

	if err == nil {
		t.Errorf("Node '%s' was not properly deleted", node.Name)
	}
}

func TestNodeStoreMemList(t *testing.T) {
	ni := NodeStoreMem{nodes: make(map[string]entities.Node)}
	node1 := entities.Node{Name: "TestRole1"}
	node2 := entities.Node{Name: "TestNode2"}

	ni.Save(node1)
	ni.Save(node1)
	ni.Save(node2)

	list := ni.List()

	if len(list) != 2 {
		t.Errorf("Nodes where not listed properly")
	}
}
