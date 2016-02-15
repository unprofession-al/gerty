// Package memstore implements the store interfaces and provides
// basic in-memory peristence.
package memstore

import (
	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/store"
)

var (
	nodes NodeStore
	roles RoleStore
)

func init() {
	store.Register("memstore", Setup)
}

// Setup allocates the required maps in memory and returns the store
// implementation. config is unused and only has the purpose to fullfil
// the func signature required by register.
func Setup(config string) (*store.Store, error) {
	nodes = NodeStore{nodes: make(map[string]entities.Node)}
	roles = RoleStore{roles: make(map[string]entities.Role)}
	s := &store.Store{
		Nodes: nodes,
		Roles: roles,
	}
	return s, nil
}
