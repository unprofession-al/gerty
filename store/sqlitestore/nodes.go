package sqlitestore

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/unprofession-al/gerty/entities"
)

var node_schema = `CREATE TABLE IF NOT EXISTS nodes (
	name text,
	vars text,
	roles text,
	PRIMARY KEY (name)
);`

// Node defines the structure of the database entries.
type Node struct {
	// name is the primary key
	Name string `db:"name"`
	// vars are stored serialized as json
	Vars string `db:"vars"`
	// roles are stored serialized as json
	Roles string `db:"roles"`
}

// NodeStore implements the entities.NodeStore interface.
type NodeStore struct {
	db *sqlx.DB
}

// Save saves/replaces a given node.
func (ns NodeStore) Save(n entities.Node) error {
	vars, err := json.Marshal(n.Vars)
	if err != nil {
		return err
	}

	roles, err := json.Marshal(n.Roles)
	if err != nil {
		return err
	}

	node := &Node{
		Name:  n.Name,
		Vars:  string(vars),
		Roles: string(roles),
	}

	_, err = ns.db.NamedExec(`INSERT OR REPLACE INTO
		nodes(name, vars, roles)
		VALUES(:name, :vars, :roles)`, node)

	return err
}

// Delete deletes a given node.
func (ns NodeStore) Delete(n entities.Node) error {
	node := &Node{Name: n.Name}

	_, err := ns.db.NamedExec(`DELETE FROM nodes
		WHERE name = :name`, node)

	return err
}

// Get retireves a node by its name.
func (ns NodeStore) Get(name string) (entities.Node, error) {
	n := Node{}

	err := ns.db.Get(&n, "SELECT * FROM nodes WHERE name=$1", name)
	if err != nil {
		return entities.Node{}, err
	}

	vars := entities.VarCollection{}
	err = json.Unmarshal([]byte(n.Vars), &vars)
	if err != nil {
		return entities.Node{}, err
	}

	roles := []string{}
	err = json.Unmarshal([]byte(n.Roles), &roles)
	if err != nil {
		return entities.Node{}, err
	}

	node := entities.Node{
		Name:  n.Name,
		Vars:  vars,
		Roles: roles,
	}

	return node, nil
}

// List returns a list of persisted nodes by their names.
func (ns NodeStore) List() ([]string, error) {
	out := []string{}

	err := ns.db.Select(&out, "SELECT name FROM nodes")

	return out, err
}
