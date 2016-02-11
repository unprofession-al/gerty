package memstore

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

type NodeStore struct {
	nodes map[string]entities.Node
}

func (ns NodeStore) Save(n entities.Node) error {
	ns.nodes[n.Name] = n
	return nil
}

func (ns NodeStore) Delete(n entities.Node) error {
	delete(ns.nodes, n.Name)
	return nil
}

func (ns NodeStore) Get(name string) (entities.Node, error) {
	node, ok := ns.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

func (ns NodeStore) List() []string {
	out := []string{}
	for name, _ := range ns.nodes {
		out = append(out, name)
	}
	return out
}
