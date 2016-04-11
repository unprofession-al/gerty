package entities

import "errors"

// Role holds information and variables as well as position information in the
// role tree
type Role struct {
	Name     string        `json:"name"`
	Vars     VarCollection `json:"vars"`
	Parent   string        `json:"parent"`
	Children []string      `json:"childern"`
}

// RoleStore needs to be implenented by the store package. It provides access
// to the persistence layer to the methods defined on the NodeInteractor.
type RoleStore interface {
	Save(r Role) error
	Delete(r Role) error
	Get(name string) (Role, error)
	HasParent(r Role) bool
	List() ([]string, error)
}

// RoleInteractor couples common actions with its persistence layer.
type RoleInteractor struct {
	RoleStore
}

// NewRoleInteractor takes the required store implementation and returns
// a RoleInteractor in order to make its methodes available.
func NewRoleInteractor(roles RoleStore) RoleInteractor {
	return RoleInteractor{RoleStore: roles}
}

// LinkParent adds a given role to the current roles children.
func (ri RoleInteractor) LinkParent(role *Role, parent *Role) error {
	if ri.checkCircularReflexion(*parent, *role) {
		return errors.New("cannot add, this would create a circular reflection")
	}
	if ri.HasParent(*role) {
		return errors.New("already has parent")
	}
	role.Parent = parent.Name
	err := ri.Save(*role)
	if err != nil {
		return err
	}
	// TODO: make sure only parents are linked on Children on Store level.
	parent.Children = append(parent.Children, role.Name)
	ri.Save(*parent)
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
// TODO: Refactor to UnlinkParent in order to be consistent with LinkParent()
func (ri RoleInteractor) UnlinkChild(role *Role, child *Role) error {
	for i, cName := range role.Children {
		c, err := ri.Get(cName)
		if err != nil {
			return err
		}
		if c.Name == child.Name {
			child.Parent = ""
			ri.Save(*child)
			// TODO: make sure only parents are linked on Children on Store level.
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
