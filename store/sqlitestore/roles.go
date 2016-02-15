package sqlitestore

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/unprofession-al/gerty/entities"
)

var roleSchema = ` CREATE TABLE IF NOT EXISTS roles (
	name text,
	vars text,
	parent text,
	children text,
	PRIMARY KEY (name)
);`

// Role defines the structure of the database entries.
type Role struct {
	Name     string `db:"name"`
	Vars     string `db:"vars"`
	Parent   string `db:"parent"`
	Children string `db:"children"`
}

// RoleStore implements the entities.RoleStore interface.
type RoleStore struct {
	db *sqlx.DB
}

// Save saves/replaces a given role.
func (rs RoleStore) Save(r entities.Role) error {
	vars, err := json.Marshal(r.Vars)
	if err != nil {
		return err
	}

	children, err := json.Marshal(r.Children)
	if err != nil {
		return err
	}

	role := &Role{
		Name:     r.Name,
		Vars:     string(vars),
		Parent:   r.Parent,
		Children: string(children),
	}

	_, err = rs.db.NamedExec(`INSERT OR REPLACE INTO
		roles(name, vars, parent, children)
		VALUES(:name, :vars, :parent, :children)`, role)

	return err
}

// Delete deletes a given role.
func (rs RoleStore) Delete(r entities.Role) error {
	role := &Role{Name: r.Name}

	_, err := rs.db.NamedExec(`DELETE FROM roles
		WHERE name = :name`, role)

	return err
}

// Get retireves a role by its name.
func (rs RoleStore) Get(name string) (entities.Role, error) {

	r := Role{}

	err := rs.db.Get(&r, "SELECT * FROM roles WHERE name=$1", name)
	if err != nil {
		return entities.Role{}, err
	}

	vars := entities.VarCollection{}
	err = json.Unmarshal([]byte(r.Vars), &vars)
	if err != nil {
		return entities.Role{}, err
	}

	children := []string{}
	err = json.Unmarshal([]byte(r.Children), &children)
	if err != nil {
		return entities.Role{}, err
	}

	role := entities.Role{
		Name:     r.Name,
		Vars:     vars,
		Parent:   r.Parent,
		Children: children,
	}

	return role, nil
}

// List returns a list of persisted roles by their names.
func (rs RoleStore) List() ([]string, error) {
	out := []string{}

	err := rs.db.Select(&out, "SELECT name FROM roles")

	return out, err
}
