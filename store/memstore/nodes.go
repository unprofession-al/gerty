package memstore

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

// NodeStore implements the entities.NodeStore interface.
type NodeStore struct {
	nodes map[string]entities.Node
}

// Save saves/replaces a given node.
func (ns NodeStore) Save(n entities.Node) error {
	ns.nodes[n.Name] = n
	return nil
}

// Delete deletes a given node.
func (ns NodeStore) Delete(n entities.Node) error {
	delete(ns.nodes, n.Name)
	return nil
}

// Get retireves a node by its name.
func (ns NodeStore) Get(name string) (entities.Node, error) {
	node, ok := ns.nodes[name]
	if !ok {
		return node, errors.New("Node does not exist")
	}
	return node, nil
}

// List returns a list of persisted nodes by their names.
func (ns NodeStore) List() ([]string, error) {
	out := []string{}
	for name := range ns.nodes {
		out = append(out, name)
	}
	return out, nil
}
