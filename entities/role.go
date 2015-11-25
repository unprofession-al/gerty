package entities

import "errors"

type Role struct {
	ID       int
	Name     string
	Vars     VarCollection
	Parent   *Role
	Children Roles
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

func (role *Role) Depth() int {
	depth := 0
	r := role
	for true {
		if r.Parent != nil {
			r = r.Parent
			depth += 1
		} else {
			break
		}
	}

	return depth
}

type Roles []*Role

func (r Roles) Len() int      { return len(r) }
func (r Roles) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Roles) Less(i, j int) bool {
	iDepth := r[i].Depth()
	jDepth := r[j].Depth()
	if iDepth == jDepth {
		return r[i].Name < r[j].Name
	}
	return iDepth > jDepth
}
