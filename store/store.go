// Package store provides a factory for the various persistence
// layers available.
package store

import (
	"errors"
	"sync"

	"github.com/unprofession-al/gerty/entities"
)

var (
	sMu sync.Mutex
	s   = make(map[string]func(string) (*Store, error))
)

// Store encapsules the NodeStore and RoleStore implementations.
type Store struct {
	Nodes entities.NodeStore
	Roles entities.RoleStore
}

// Register is called in the init() funcs of the actual Store
// implemetations in order to make the implementation available
// to the factory.
func Register(name string, setupFunc func(string) (*Store, error)) {
	sMu.Lock()
	defer sMu.Unlock()
	if _, dup := s[name]; dup {
		panic("store: Register called twice for store " + name)
	}
	s[name] = setupFunc
}

// New returns a Store containg the sperified and configured NodeStore and
// RoleStore implemetations.
func New(name string, config string) (*Store, error) {
	setupFunc, ok := s[name]
	if !ok {
		return nil, errors.New("store: store '" + name + "' does not exist")
	}
	return setupFunc(config)
}
