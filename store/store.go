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

type Store struct {
	Nodes entities.NodeStore
	Roles entities.RoleStore
}

func Register(name string, setupFunc func(string) (*Store, error)) {
	sMu.Lock()
	defer sMu.Unlock()
	if _, dup := s[name]; dup {
		panic("store: Register called twice for store " + name)
	}
	s[name] = setupFunc
}

func New(name string, config string) (*Store, error) {
	setupFunc, ok := s[name]
	if !ok {
		return nil, errors.New("store: store '" + name + "' does not exist")
	}
	return setupFunc(config)
}
