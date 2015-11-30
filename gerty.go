package main

import (
	"fmt"

	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/store"
)

func main() {
	roleMockData := map[int64]*store.RoleMock{
		1: {Name: "ra", Children: []int64{2, 3}},
		2: {Name: "rb"},
		3: {Name: "rc"},
	}
	roleBackendMock := store.NewRoleBackendMock(roleMockData)

	nodeMockData := []*store.NodeMock{
		{Name: "Node 1", Roles: []int64{2, 3}},
	}
	nodeBackendMock := store.NewNodeBackendMock(nodeMockData, &roleBackendMock)

	store := store.Store{
		Roles: &roleBackendMock,
		Nodes: &nodeBackendMock,
	}

	role, _ := store.Roles.Get("rc")
	roles := entities.Roles{}
	_ = append(roles, &role)

	fmt.Printf("Nothing to see here yet, move along.\n")
}
