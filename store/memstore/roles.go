package memstore

import (
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

type RoleStore struct {
	roles map[string]entities.Role
}

func (rs RoleStore) Save(r entities.Role) error {
	rs.roles[r.Name] = r
	return nil
}

func (rs RoleStore) Delete(r entities.Role) error {
	delete(rs.roles, r.Name)
	return nil
}

func (rs RoleStore) Get(name string) (entities.Role, error) {
	role, ok := rs.roles[name]
	if !ok {
		return role, errors.New("Role does not exist")
	}
	return role, nil
}

func (rs RoleStore) List() ([]string, error) {
	out := []string{}
	for name, _ := range rs.roles {
		out = append(out, name)
	}
	return out, nil
}
