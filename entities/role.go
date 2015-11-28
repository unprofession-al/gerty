package entities

import "errors"

// Role holds information and variables as well as position information in the
// role tree
type Role struct {
	Name     string
	Vars     VarCollection
	Parent   *Role
	Children Roles
}

// LinkChild adds a given role to the current roles children.
func (role *Role) LinkChild(child *Role) error {
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
		if c == role {
			return true
		}
		return role.checkCircularReflexion(c)
	}
	return false
}

// UnlinkChild removes a given role from the current roles children.
func (role *Role) UnlinkChild(child *Role) error {
	for i, c := range role.Children {
		if c == child {
			child.Parent = &Role{}
			role.Children = append(role.Children[:i], role.Children[i+1:]...)
			return nil
		}
	}
	return errors.New("child is not related to this parent")
}

// Depth calculates the distance of the role to the root of the tree.
func (role *Role) Depth() int {
	depth := 0
	r := role
	for true {
		if r.GetParent() != nil {
			r = r.GetParent()
			depth++
		} else {
			break
		}
	}

	return depth
}

// GetParent returns the parent role
func (role *Role) GetParent() *Role {
	if role.Parent != nil {
		return role.Parent
	}
	return nil
}

// Roles is a list of references to `Role` elements
type Roles []*Role

// Len returns the lenght of the list, part of implementing the sort interface.
func (r Roles) Len() int { return len(r) }

// Swap changes the order of two elements in the list, part of implementing
// the sort interface.
func (r Roles) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Less compares the order of two elements, part of implementing the sort interface.
// The heigher the depht (distance to root element of the tree), the earlier the
// element appears in the the order. Roles with the same depth will be sorted
// alphabetically.
func (r Roles) Less(i, j int) bool {
	iDepth := r[i].Depth()
	jDepth := r[j].Depth()
	if iDepth == jDepth {
		return r[i].Name < r[j].Name
	}
	return iDepth > jDepth
}
