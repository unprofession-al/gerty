package memstore

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

// RoleStore implements the entities.RoleStore interface.
type RoleStore struct {
	roles map[string]entities.Role
}

// Save saves/replaces a given role.
func (rs RoleStore) Save(r entities.Role) error {
	rs.roles[r.Name] = r
	return nil
}

// Delete deletes a given role.
func (rs RoleStore) Delete(r entities.Role) error {
	delete(rs.roles, r.Name)
	return nil
}

// Get retireves a role by its name.
func (rs RoleStore) Get(name string) (entities.Role, error) {
	role, ok := rs.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}

// List returns a list of persisted roles by their names.
func (rs RoleStore) List() ([]string, error) {
	out := []string{}
	for name := range rs.roles {
		out = append(out, name)
	}
	return out, nil
}
