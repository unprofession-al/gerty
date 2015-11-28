package store

import "github.com/unprofession-al/gerty/entities"

type Store struct {
	Nodes NodeBackend
	Roles RoleBackend
}

type NodeBackend interface {
	Save(n entities.Node) error
	Delete(n entities.Node) error
	Get(name string) (entities.Node, error)
}

type RoleBackend interface {
	Save(n entities.Role) error
	Delete(n entities.Role) error
	Get(name string) (entities.Role, error)
}
