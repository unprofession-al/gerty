package entities

import "errors"

type RoleStore interface {
	Save(r Role) error
	Delete(r Role) error
	Get(name string) (Role, error)
	List() ([]string, error)
}

type RoleInteractor struct {
	RoleStore
}

func NewRoleInteractor(roles RoleStore) RoleInteractor {
	return RoleInteractor{RoleStore: roles}
}

// Role holds information and variables as well as position information in the
// role tree
type Role struct {
	Name     string        `json:"name"`
	Vars     VarCollection `json:"vars"`
	Parent   string        `json:"parent"`
	Children []string      `json:"childern"`
}

// LinkChild adds a given role to the current roles children.
func (ri RoleInteractor) LinkChild(role *Role, child *Role) error {
	if ri.checkCircularReflexion(*role, *child) {
		return errors.New("cannot add, this would create a circular reflection")
	}
	if child.Parent != "" {
		return errors.New("child already has parent")
	}
	role.Children = append(role.Children, child.Name)
	ri.Save(*role)
	child.Parent = role.Name
	ri.Save(*child)
	return nil
}

func (ri RoleInteractor) checkCircularReflexion(role Role, child Role) bool {
	for _, cName := range child.Children {
		c, err := ri.Get(cName)
		if err != nil {
			panic(err)
		}
		if c.Name == role.Name {
			return true
		}
		return ri.checkCircularReflexion(role, c)
	}
	return false
}

// UnlinkChild removes a given role from the current roles children.
func (ri RoleInteractor) UnlinkChild(role *Role, child *Role) error {
	for i, cName := range role.Children {
		c, err := ri.Get(cName)
		if err != nil {
			return err
		}
		if c.Name == child.Name {
			child.Parent = ""
			ri.Save(*child)
			role.Children = append(role.Children[:i], role.Children[i+1:]...)
			ri.Save(*role)
			return nil
		}
	}
	return errors.New("child is not related to this parent")
}

// Depth calculates the distance of the role to the root of the tree.
func (ri RoleInteractor) Depth(role Role) int {
	depth := 0
	r := role
	for true {
		if r.Parent != "" {
			nr, err := ri.Get(r.Parent)
			if err != nil {
				panic(err)
			}
			r = nr
			depth++
		} else {
			break
		}
	}
	return depth
}

// Roles is a list of references to `Role` elements
type roleSorter struct {
	Roles []string
	ri    RoleInteractor
}

// Len returns the length of the list, part of implementing the sort interface.
func (r roleSorter) Len() int { return len(r.Roles) }

// Swap changes the order of two elements in the list, part of implementing
// the sort interface.
func (r roleSorter) Swap(i, j int) { r.Roles[i], r.Roles[j] = r.Roles[j], r.Roles[i] }

// Less compares the order of two elements, part of implementing the sort interface.
// The higher the depht (distance to root element of the tree), the earlier the
// element appears in the the order. Roles with the same depth will be sorted
// alphabetically.
func (r roleSorter) Less(i, j int) bool {
	iRole, err := r.ri.Get(r.Roles[i])
	if err != nil {
		panic(err)
	}
	iDepth := r.ri.Depth(iRole)

	jRole, err := r.ri.Get(r.Roles[j])
	if err != nil {
		panic(err)
	}
	jDepth := r.ri.Depth(jRole)

	if iDepth == jDepth {
		return iRole.Name < jRole.Name
	}
	return iDepth > jDepth
}
