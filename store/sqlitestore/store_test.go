package sqlitestore

import (
	"os"
	"testing"

	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/store"
)

const filename = "/tmp/sqlitestore_test.sqlite3"

func TestStore(t *testing.T) {
	// init store
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			t.Errorf(err.Error())
		}
	}
	s, err := store.New("sqlitestore", filename)
	if err != nil {
		t.Errorf("Store could not be connected: %s", err.Error())
	}

	// create role
	roleOrig := entities.Role{Name: "testrole"}
	if err := s.Roles.Save(roleOrig); err != nil {
		t.Errorf("Role could not be stored: %s", err.Error())
	}

	// create role with parent
	childOrig := entities.Role{Name: "childrole", Parent: roleOrig.Name}
	if err := s.Roles.Save(childOrig); err != nil {
		t.Errorf("Role could not be stored: %s", err.Error())
	}

	// try to delete parent, which should fail
	if err := s.Roles.Delete(roleOrig); err == nil {
		t.Error("Parent could be deleted, should not")
	}

	// create another role with parent
	grandchildOrig := entities.Role{Name: "grandchildrole", Parent: childOrig.Name}
	if err := s.Roles.Save(grandchildOrig); err != nil {
		t.Errorf("Role could not be stored: %s", err.Error())
	}

	// fetch child
	_, err = s.Roles.Get("childrole")
	if err != nil {
		t.Errorf("Role could not be fetched: %s", err.Error())
	}

	// fetch roles
	if _, err := s.Roles.List(); err != nil {
		t.Errorf("Could not fetch role list, %s", err.Error())
	}

	// save node with roles
	nodeOrig := entities.Node{Name: "testnode", Roles: []string{childOrig.Name, grandchildOrig.Name}}
	if err := s.Nodes.Save(nodeOrig); err != nil {
		t.Errorf("Node could not be stored: %s", err.Error())
	}

	// delete grandchild
	if err := s.Roles.Delete(grandchildOrig); err != nil {
		t.Errorf("Grandchild could not be deleted: %s", err.Error())
	}

	// fetch node
	nodeFetched, err := s.Nodes.Get("testnode")
	if err != nil {
		t.Errorf("Node could not be fetched: %s", err.Error())
	}

	// inspert node
	if roleCount := len(nodeFetched.Roles); roleCount != 1 {
		t.Errorf("Node has wrong number of Roles: %d, should be 1", roleCount)
	}

	// delete child
	if err := s.Roles.Delete(childOrig); err != nil {
		t.Errorf("Child could not be deleted: %s", err.Error())
	}

	// fetch node
	nodeFetched, err = s.Nodes.Get("testnode")
	if err != nil {
		t.Errorf("Node could not be fetched: %s", err.Error())
	}

	// inspert node
	if roleCount := len(nodeFetched.Roles); roleCount != 0 {
		t.Errorf("Node has wrong number of Roles: %d, should be 1", roleCount)
	}

	// delete child

	// fetch nodes
	if _, err := s.Nodes.List(); err != nil {
		t.Errorf("Could not fetch node list, %s", err.Error())
	}

	// delete node
	if err := s.Nodes.Delete(nodeFetched); err != nil {
		t.Errorf("Node could not be deleted: %s", err.Error())
	}
}
