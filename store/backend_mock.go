package store

import (
	"encoding/json"
	"errors"

	"github.com/unprofession-al/gerty/entities"
)

type RoleMock struct {
	Name     string
	Vars     string
	Children []int64
}

type RoleBackendMock map[int64]*entities.Role

func NewRoleBackendMock(mockData map[int64]*RoleMock) RoleBackendMock {
	mocked := RoleBackendMock{}
	for id, roleMock := range mockData {
		if id == 0 {
			panic("cannot init role mock backend if id is `0`")
		}
		var vars entities.VarCollection

		if roleMock.Vars == "" {
			roleMock.Vars = "[]"
		}
		err := json.Unmarshal([]byte(roleMock.Vars), &vars)
		if err != nil {
			panic("JSON not valid: " + roleMock.Vars)
		}

		mocked[id] = &entities.Role{Name: roleMock.Name, Vars: vars}
	}

	for roleId, role := range mocked {
		for _, childId := range mockData[roleId].Children {
			role.LinkChild(mocked[childId])
		}
	}
	return mocked
}

func (rm *RoleBackendMock) Save(r entities.Role) error {
	found := false
	var highest int64
	for id, role := range *rm {
		if role.Name == r.Name {
			(*rm)[id] = &r
			found = true
			break
		}
		if id > highest {
			highest = id
		}
	}
	if !found {
		(*rm)[highest+1] = &r
	}
	return nil
}

func (rm *RoleBackendMock) Delete(r entities.Role) error {
	for id, role := range *rm {
		if role.Name == r.Name {
			delete(*rm, id)
			return nil
		}
	}
	return errors.New("Role does not exist and therefore cannot be deleted")
}

func (rm RoleBackendMock) Get(name string) (entities.Role, error) {
	for _, role := range rm {
		if role.Name == name {
			return *role, nil
		}
	}
	return entities.Role{}, errors.New("Role not found")
}

func (rm RoleBackendMock) getId(name string) int64 {
	for id, role := range rm {
		if role.Name == name {
			return id
		}
	}
	return 0
}

type NodeMock struct {
	Name  string
	Vars  string
	Roles []int64
}

type NodeBackendMock struct {
	nodes []*NodeMock
	roles *RoleBackendMock
}

func NewNodeBackendMock(mockData []*NodeMock, roles *RoleBackendMock) NodeBackendMock {
	mocked := NodeBackendMock{
		nodes: mockData,
		roles: roles,
	}
	return mocked
}

func (nm *NodeBackendMock) Save(n entities.Node) error {
	b, err := json.Marshal(n.Vars)
	if err != nil {
		return errors.New("could not marshal node vars: " + err.Error())
	}
	vars := string(b)

	var roles []int64
	for _, role := range n.Roles {
		roleId := nm.roles.getId(role.Name)
		if roleId != 0 {
			roles = append(roles, roleId)
		}
	}

	nodeMock := &NodeMock{
		Name:  n.Name,
		Vars:  vars,
		Roles: roles,
	}

	for id, node := range nm.nodes {
		if nodeMock.Name == node.Name {
			nm.nodes[id] = nodeMock
			return nil
		}
	}

	nm.nodes = append(nm.nodes, nodeMock)
	return nil
}

func (nm *NodeBackendMock) Delete(n entities.Node) error {
	for id, node := range nm.nodes {
		if n.Name == node.Name {
			nm.nodes = append(nm.nodes[:id], nm.nodes[id+1:]...)
			return nil
		}
	}
	return errors.New("node not found")
}

func (nm NodeBackendMock) Get(name string) (entities.Node, error) {
	for _, node := range nm.nodes {
		if name == node.Name {
			var vars entities.VarCollection
			err := json.Unmarshal([]byte(node.Vars), &vars)
			if err != nil {
				return entities.Node{}, errors.New("JSON not valid")
			}

			roles := entities.Roles{}

			for _, roleId := range node.Roles {
				roles = append(roles, (*nm.roles)[roleId])
			}

			return entities.Node{
				Name:  name,
				Vars:  vars,
				Roles: roles,
			}, nil
		}
	}
	return entities.Node{}, errors.New("node not found")
}
