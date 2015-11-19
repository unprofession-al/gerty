package main

import "errors"

type Role struct {
	ID       int
	Name     string
	Vars     VarCollection
	Parent   *Role
	Children []*Role
}

func (role *Role) AddChild(child *Role) error {
	if role.checkCircularReflexion(child) {
		return errors.New("cannot add, this would create a circular reflection")
	}
	if child.Parent != nil {
		return errors.New("child already has parent")
	}
	role.Children = append(role.Children, child)
	child.Parent = role
	return nil
}

func (role *Role) checkCircularReflexion(child *Role) bool {
	for _, c := range child.Children {
		if c.ID == role.ID {
			return true
		}
		return role.checkCircularReflexion(c)
	}
	return false
}

func (role *Role) DeleteChild(child *Role) error {
	for i, c := range role.Children {
		if c.ID == child.ID {
			child.Parent = &Role{}
			role.Children = append(role.Children[:i], role.Children[i+1:]...)
			return nil
		}
	}
	return errors.New("child is not related to this parent")
}
