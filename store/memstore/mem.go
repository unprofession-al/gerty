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
	store.Register("mem", Setup)
}

func Setup(config string) (*store.Store, error) {
	nodes = NodeStore{nodes: make(map[string]entities.Node)}
	roles = RoleStore{roles: make(map[string]entities.Role)}
	s := &store.Store{
		Nodes: nodes,
		Roles: roles,
	}
	return s, nil
}
