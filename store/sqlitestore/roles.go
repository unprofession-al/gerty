package sqlitestore

import (
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/unprofession-al/gerty/entities"
)

// RoleStore implements the entities.RoleStore interface.
type RoleStore struct {
	db *sqlx.DB
}

// Save saves/replaces a given role.
func (rs RoleStore) Save(r entities.Role) error {
	role := &Role{Name: r.Name}

	// serialize vars
	vars, err := json.Marshal(r.Vars)
	if err != nil {
		return err
	}
	role.Vars = string(vars)

	// setup transaction
	tx, err := rs.db.Beginx()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// get id of parent role
	var parentID sql.NullInt64
	if r.Parent != "" {
		err = tx.Get(&parentID, "SELECT id FROM role WHERE name = $1;", r.Parent)
		if err != nil {
			return err
		}
		role.Parent = parentID
	}

	_, err = tx.NamedExec(`INSERT OR REPLACE INTO
		role(name, vars, parent)
		VALUES(:name, :vars, :parent);`, role)

	return err
}

// Delete deletes a given role.
func (rs RoleStore) Delete(r entities.Role) error {
	role := &Role{Name: r.Name}

	_, err := rs.db.NamedExec(`DELETE FROM role
		WHERE name = :name;`, role)

	return err
}

// Get retireves a role by its name.
func (rs RoleStore) Get(name string) (entities.Role, error) {

	r := Role{}

	err := rs.db.Get(&r, `SELECT r1.id, r1.name, r1.vars, r2.name AS parent_name
	                        FROM role r1
	             LEFT OUTER JOIN role r2
	                          ON r1.parent = r2.id
	                       WHERE r1.name = $1;`, name)
	if err != nil {
		return entities.Role{}, err
	}

	vars := entities.VarCollection{}
	err = json.Unmarshal([]byte(r.Vars), &vars)
	if err != nil {
		return entities.Role{}, err
	}

	// SELECT '[\"' || GROUP_CONCAT(name,'\",\"') || '\"]' AS aoeua FROM role WHERE parent=1;
	children := []string{}
	err = rs.db.Select(&children, "SELECT name FROM role WHERE parent = $1;", r.ID)
	if err != nil {
		return entities.Role{}, err
	}

	role := entities.Role{
		Name:     r.Name,
		Vars:     vars,
		Parent:   r.ParentName,
		Children: children,
	}

	return role, nil
}

// List returns a list of persisted roles by their names.
func (rs RoleStore) List() ([]string, error) {
	out := []string{}

	err := rs.db.Select(&out, "SELECT name FROM role;")

	return out, err
}
