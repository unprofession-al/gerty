package entities

import (
	"sort"
	"testing"
)

func TestRoleInteractor(t *testing.T) {
	newRole := Role{Name: "TestRole"}

	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	ri.Save(newRole)

	role, _ := ri.Get(newRole.Name)

	if role.Name != newRole.Name {
		t.Errorf("Name in not consistent: `%s` != '%s'", role.Name, newRole.Name)
	}
}

func TestRoleAddChildToParent(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p := Role{Name: "Parent"}
	ri.Save(p)

	c := Role{Name: "Child"}
	ri.Save(c)

	err := ri.LinkParent(&c, &p)
	if err != nil {
		t.Error("child not added to parent")
	}

	p, _ = ri.Get(p.Name)
	c, _ = ri.Get(c.Name)

	childrenAdded := false
	for _, elem := range p.Children {
		if elem == c.Name {
			childrenAdded = true
		}
	}
	if !childrenAdded {
		t.Error("child not added to parent")
	}

	if c.Parent != p.Name {
		t.Error("group not refered in child")

	}
}

func TestRoleCreateCircle(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p := Role{Name: "Parent"}
	ri.Save(p)

	c1 := Role{Name: "Child Level 1"}
	ri.Save(c1)

	c2 := Role{Name: "Child Level 2"}
	ri.Save(c2)

	ri.LinkParent(&c1, &p)
	ri.LinkParent(&c2, &c1)

	err := ri.LinkParent(&p, &c2)
	if err == nil {
		t.Error("circular reflecton created")
	}
}

func TestRoleMultipleParents(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p1 := Role{Name: "Parent 1"}
	ri.Save(p1)

	p2 := Role{Name: "Parent 2"}
	ri.Save(p2)

	c := Role{Name: "Child 1"}
	ri.Save(c)

	ri.LinkParent(&c, &p1)

	err := ri.LinkParent(&c, &p2)
	if err == nil {
		t.Error("child with parent added to new parent")
	}
}

func TestRoleDeleteChildFromParent(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p := Role{Name: "Parent"}
	ri.Save(p)

	c := Role{Name: "Child"}
	ri.Save(c)

	ri.LinkParent(&c, &p)

	err := ri.UnlinkChild(&p, &c)
	if err != nil {
		t.Error("child not deleted from parent")
	}

	for _, elem := range p.Children {
		if elem == c.Name {
			t.Error("child not removed from parent")
		}
	}

	if c.Parent == p.Name {
		t.Error("child stil refers to parent")

	}
}

func TestRoleDeleteUnrelatedChild(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p := Role{Name: "Parent"}
	ri.Save(p)

	c1 := Role{Name: "Child Level 1"}
	ri.Save(c1)

	c2 := Role{Name: "Child Level 2"}
	ri.Save(c2)

	ri.LinkParent(&c1, &p)
	err := ri.UnlinkChild(&p, &c2)
	if err == nil {
		t.Error("removing a node from another node that is not related as parent shuld be impossible")
	}
}

func TestRoleCalculateDepth(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p := Role{Name: "Parent"}
	ri.Save(p)

	c1 := Role{Name: "Child Level 1"}
	ri.Save(c1)

	c2 := Role{Name: "Child Level 2"}
	ri.Save(c2)

	ri.LinkParent(&c1, &p)
	ri.LinkParent(&c2, &c1)

	if ri.Depth(c2) != 2 {
		t.Error("depth not calculated correctly")
	}
}

func TestRoleSort(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})

	p1 := Role{Name: "C Parent 1"}
	ri.Save(p1)
	c11 := Role{Name: "C Child 1 Level 1"}
	ri.Save(c11)
	c12 := Role{Name: "C Child 1 Level 2"}
	ri.Save(c12)
	ri.LinkParent(&c11, &p1)
	ri.LinkParent(&c12, &c11)

	p2 := Role{Name: "B Parent 2"}
	ri.Save(p2)
	c21 := Role{Name: "B Child 2 Level 1"}
	ri.Save(c21)
	ri.LinkParent(&c21, &p2)

	p3 := Role{Name: "A Parent 3"}
	ri.Save(p3)
	c31 := Role{Name: "A Child 3 Level 1"}
	ri.Save(c31)
	c32 := Role{Name: "A Child 3 Level 2"}
	ri.Save(c32)
	ri.LinkParent(&c31, &p3)
	ri.LinkParent(&c32, &c31)

	expected := map[int]string{
		0: "A Child 3 Level 2",
		1: "C Child 1 Level 2",
		2: "B Child 2 Level 1",
	}

	roles := []string{c12.Name, c21.Name, c32.Name}
	rs := roleSorter{
		Roles: roles,
		ri:    ri,
	}

	sort.Sort(rs)

	for k, v := range rs.Roles {
		if expected[k] != v {
			t.Errorf("`%s` is at position %d, should be `%s`", v, k, expected[k])
		}
	}
}
