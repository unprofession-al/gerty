package sqlitestore

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/unprofession-al/gerty/entities"
)

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

	// setup transaction
	tx, err := ns.db.Beginx()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// save node
	node := &Node{
		Name: n.Name,
		Vars: string(vars),
	}

	result, err := tx.NamedExec(`INSERT OR REPLACE INTO
		node(name, vars)
		VALUES(:name, :vars);`, node)

	// save relations to roles if required
	if len(n.Roles) > 0 {
		nodeID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		data := struct {
			NodeID int64    `db:"nid"`
			Roles  []string `db:"roles"`
		}{
			NodeID: nodeID,
			Roles:  n.Roles,
		}

		query, args, err := sqlx.Named(`INSERT OR IGNORE INTO
								   node_role(role_id, node_id, id)
				            SELECT r.id AS role_id,
				                   :nid AS node_id,
			                       :nid || '-' || r.ID AS id
		                      FROM role r
							  WHERE r.name IN (:roles);`, data)
		if err != nil {
			return err
		}

		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return err
		}

		query = tx.Rebind(query)
		_, err = tx.Exec(query, args...)
	}

	return err
}

// Delete deletes a given node.
func (ns NodeStore) Delete(n entities.Node) error {
	node := &Node{Name: n.Name}

	_, err := ns.db.NamedExec(`DELETE FROM node
		WHERE name = :name;`, node)

	return err
}

// Get retireves a node by its name.
func (ns NodeStore) Get(name string) (entities.Node, error) {
	n := Node{}

	err := ns.db.Get(&n, "SELECT * FROM node WHERE name=$1;", name)
	if err != nil {
		return entities.Node{}, err
	}

	vars := entities.VarCollection{}
	err = json.Unmarshal([]byte(n.Vars), &vars)
	if err != nil {
		return entities.Node{}, err
	}

	roles := []string{}
	err = ns.db.Select(&roles, `SELECT r.name
		                          FROM node_role nr
		               LEFT OUTER JOIN role r
			      				    ON nr.role_id = r.id
		                         WHERE nr.node_id = $1;`, n.ID)
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

	err := ns.db.Select(&out, "SELECT name FROM node;")

	return out, err
}
