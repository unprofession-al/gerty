package store

import (
	"errors"
	"sync"

	"github.com/unprofession-al/gerty/entities"
)

var (
	stMu sync.Mutex
	st   = make(map[string]Stores)
)

type Stores struct {
	Nodes entities.NodeStore
	Roles entities.RoleStore
}

func register(name string, nodes entities.NodeStore, roles entities.RoleStore) {
	stMu.Lock()
	defer stMu.Unlock()
	if nodes == nil || roles == nil {
		panic("store: Register store is nil")
	}
	if _, dup := st[name]; dup {
		panic("store: Register called twice for store " + name)
	}
	st[name] = Stores{
		Nodes: nodes,
		Roles: roles,
	}
}

func Open(name string, config string) (Stores, error) {
	stores, ok := st[name]
	if !ok {
		return stores, errors.New("store: store '" + name + "' does not exist")
	}
	return st[name], nil
}
