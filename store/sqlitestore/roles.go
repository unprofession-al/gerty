package sqlitestore

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/unprofession-al/gerty/entities"
)

const rootRoleName = "all"

// RoleStore implements the entities.RoleStore interface.
type RoleStore struct {
	db *sqlx.DB
}

// Save saves/replaces a given role.
func (rs RoleStore) Save(r entities.Role) error {
	role := &Role{Name: r.Name}

	// serialize vars
	vars, err := r.Vars.Serialize()
	if err != nil {
		return err
	}
	role.Vars = string(vars)

	if r.Parent != "" && r.Name == rootRoleName {
		fmt.Println("eee")
		return errors.New("Cannot add parent role to root role")
	}

	// get id of parent role
	var parentID sql.NullInt64
	if r.Parent != "" {
		err = rs.db.Get(&parentID, "SELECT id FROM role WHERE name = $1;", r.Parent)
		if err != nil {
			return err
		}
		role.Parent = parentID
	}

	// figure out if role already exists
	err = rs.db.Get(&role.ID, "SELECT id FROM role WHERE name = $1;", r.Name)
	if err == nil {
		_, err = rs.db.NamedExec(`UPDATE role
			SET name = :name,
				vars = :vars,
				parent = :parent
		  WHERE id = :id;`, role)
	} else {
		_, err = rs.db.NamedExec(`INSERT INTO
			role(name, vars, parent)
			VALUES(:name, :vars, :parent);`, role)
	}

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

	err = rs.db.Select(&r.Hosts, `SELECT n.name
		                          FROM node_role nr
		               LEFT OUTER JOIN node n
			      				    ON nr.node_id = n.id
		                         WHERE nr.role_id = $1;`, r.ID)
	if err != nil {
		return entities.Role{}, err
	}

	// SELECT '[\"' || GROUP_CONCAT(name,'\",\"') || '\"]' AS aoeua FROM role WHERE parent=1;
	err = rs.db.Select(&r.Children, "SELECT name FROM role WHERE parent = $1;", r.ID)
	if err != nil {
		return entities.Role{}, err
	}

	err = r.populateRootRole(rs)
	if err != nil {
		return entities.Role{}, err
	}

	role := entities.Role{
		Name:     r.Name,
		Vars:     entities.VarCollection{},
		Parent:   r.ParentName.String,
		Children: r.Children,
		Nodes:    r.Hosts,
	}

	err = role.Vars.Deserialize([]byte(r.Vars))
	if err != nil {
		return entities.Role{}, err
	}

	return role, nil
}

func (r *Role) populateRootRole(rs RoleStore) error {
	// Check if rootRole exisits in database
	var rr int
	err := rs.db.Get(&rr, "SELECT id FROM role WHERE name = $1;", rootRoleName)
	if err != nil {
		return nil
	}

	if r.ParentName.Valid == false && r.Name != rootRoleName {
		// add all as parent
		r.ParentName.String = rootRoleName
		r.ParentName.Valid = true
	} else if r.Name == rootRoleName {
		// add all roles without parent as children
		children := []string{}
		err = rs.db.Select(&children, "SELECT name FROM role WHERE parent IS NULL AND id != $1", rr)
		if err != nil {
			return err
		}
		r.Children = append(r.Children, children...)
	}
	return nil
}

// List returns a list of persisted roles by their names.
func (rs RoleStore) List() ([]string, error) {
	out := []string{}

	err := rs.db.Select(&out, "SELECT name FROM role;")

	return out, err
}

// Delete deletes a given role.
func (rs RoleStore) HasParent(r entities.Role) bool {
	if r.Parent == rootRoleName {
		return false
	} else {
		return true
	}
}
