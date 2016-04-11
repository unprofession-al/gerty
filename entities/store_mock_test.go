package entities

import "errors"

type NodeStoreMock struct {
	nodes map[string]Node
}

func (ns NodeStoreMock) Save(n Node) error {
	ns.nodes[n.Name] = n
	return nil
}

func (ns NodeStoreMock) Delete(n Node) error {
	delete(ns.nodes, n.Name)
	return nil
}

func (ns NodeStoreMock) Get(name string) (Node, error) {
	node, ok := ns.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

func (ns NodeStoreMock) List() ([]string, error) {
	out := []string{}
	for name := range ns.nodes {
		out = append(out, name)
	}
	return out, nil
}

type RoleStoreMock struct {
	roles map[string]Role
}

func (rs RoleStoreMock) Save(r Role) error {
	rs.roles[r.Name] = r
	return nil
}

func (rs RoleStoreMock) Delete(r Role) error {
	delete(rs.roles, r.Name)
	return nil
}

func (rs RoleStoreMock) Get(name string) (Role, error) {
	role, ok := rs.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}

func (rs RoleStoreMock) HasParent(r Role) bool {
	if _, ok := rs.roles[r.Parent]; ok {
		return true
	} else {
		return false
	}
}

func (rs RoleStoreMock) List() ([]string, error) {
	out := []string{}
	for name := range rs.roles {
		out = append(out, name)
	}
	return out, nil
}
