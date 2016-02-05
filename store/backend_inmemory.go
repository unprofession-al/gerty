package store

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

type NodeStoreImpl struct {
	nodes map[string]entities.Node
}

func (nsi NodeStoreImpl) Save(n entities.Node) error {
	nsi.nodes[n.Name] = n
	return nil
}

func (nsi NodeStoreImpl) Delete(n entities.Node) error {
	delete(nsi.nodes, n.Name)
	return nil
}

func (nsi NodeStoreImpl) Get(name string) (entities.Node, error) {
	node, ok := nsi.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

type RoleStoreImpl struct {
	roles map[string]entities.Role
}

func (rsi RoleStoreImpl) Save(r entities.Role) error {
	rsi.roles[r.Name] = r
	return nil
}

func (rsi RoleStoreImpl) Delete(r entities.Role) error {
	delete(rsi.roles, r.Name)
	return nil
}

func (rsi RoleStoreImpl) Get(name string) (entities.Role, error) {
	role, ok := rsi.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}
