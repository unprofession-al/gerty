package store

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

func init() {
	nodes := NodeStoreMem{nodes: make(map[string]entities.Node)}
	roles := RoleStoreMem{roles: make(map[string]entities.Role)}
	register("mem", nodes, roles)
}

type NodeStoreMem struct {
	nodes map[string]entities.Node
}

func (ns NodeStoreMem) Save(n entities.Node) error {
	ns.nodes[n.Name] = n
	return nil
}

func (ns NodeStoreMem) Delete(n entities.Node) error {
	delete(ns.nodes, n.Name)
	return nil
}

func (ns NodeStoreMem) Get(name string) (entities.Node, error) {
	node, ok := ns.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

func (ns NodeStoreMem) List() []string {
	out := []string{}
	for name, _ := range ns.nodes {
		out = append(out, name)
	}
	return out
}

type RoleStoreMem struct {
	roles map[string]entities.Role
}

func (rs RoleStoreMem) Save(r entities.Role) error {
	rs.roles[r.Name] = r
	return nil
}

func (rs RoleStoreMem) Delete(r entities.Role) error {
	delete(rs.roles, r.Name)
	return nil
}

func (rs RoleStoreMem) Get(name string) (entities.Role, error) {
	role, ok := rs.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}

func (rs RoleStoreMem) List() []string {
	out := []string{}
	for name, _ := range rs.roles {
		out = append(out, name)
	}
	return out
}