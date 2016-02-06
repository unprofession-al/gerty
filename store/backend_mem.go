package store

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

type NodeStore struct {
	nodes map[string]entities.Node
}

func (nsi NodeStore) Save(n entities.Node) error {
	nsi.nodes[n.Name] = n
	return nil
}

func (nsi NodeStore) Delete(n entities.Node) error {
	delete(nsi.nodes, n.Name)
	return nil
}

func (nsi NodeStore) Get(name string) (entities.Node, error) {
	node, ok := nsi.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

type RoleStore struct {
	roles map[string]entities.Role
}

func (rsi RoleStore) Save(r entities.Role) error {
	rsi.roles[r.Name] = r
	return nil
}

func (rsi RoleStore) Delete(r entities.Role) error {
	delete(rsi.roles, r.Name)
	return nil
}

func (rsi RoleStore) Get(name string) (entities.Role, error) {
	role, ok := rsi.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}
