package entities

import "errors"

type NodeStoreMock struct {
	nodes map[string]Node
}

func (nsi NodeStoreMock) Save(n Node) error {
	nsi.nodes[n.Name] = n
	return nil
}

func (nsi NodeStoreMock) Delete(n Node) error {
	delete(nsi.nodes, n.Name)
	return nil
}

func (nsi NodeStoreMock) Get(name string) (Node, error) {
	node, ok := nsi.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

type RoleStoreMock struct {
	roles map[string]Role
}

func (rsi RoleStoreMock) Save(r Role) error {
	rsi.roles[r.Name] = r
	return nil
}

func (rsi RoleStoreMock) Delete(r Role) error {
	delete(rsi.roles, r.Name)
	return nil
}

func (rsi RoleStoreMock) Get(name string) (Role, error) {
	role, ok := rsi.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}
