package entities

import (
	"sort"
	"testing"
)

func TestRoleAddChildToParent(t *testing.T) {
	p := Role{Name: "Parent"}
	c := Role{Name: "Child"}

	err := p.LinkChild(&c)
	if err != nil {
		t.Error("child not added to parent")
	}

	childrenAdded := false
	for _, elem := range p.Children {
		if elem == &c {
			childrenAdded = true
		}
	}
	if !childrenAdded {
		t.Error("child not added to parent")
	}

	if c.Parent != &p {
		t.Error("group not refered in child")

	}
}

func TestRoleCalculateDepth(t *testing.T) {
	p := Role{Name: "Parent"}
	c1 := Role{Name: "Child Level 1"}
	c2 := Role{Name: "Child Level 2"}

	p.LinkChild(&c1)
	c1.LinkChild(&c2)

	if c2.Depth() != 2 {
		t.Error("depth not calculated correctly")
	}
}

func TestRoleSort(t *testing.T) {
	p1 := Role{Name: "C Parent 1"}
	c11 := Role{Name: "C Child 1 Level 1"}
	c12 := Role{Name: "C Child 1 Level 2"}
	p1.LinkChild(&c11)
	c11.LinkChild(&c12)

	p2 := Role{Name: "B Parent 2"}
	c21 := Role{Name: "B Child 2 Level 1"}
	p2.LinkChild(&c21)

	p3 := Role{Name: "A Parent 3"}
	c31 := Role{Name: "A Child 3 Level 1"}
	c32 := Role{Name: "A Child 3 Level 2"}
	p3.LinkChild(&c31)
	c31.LinkChild(&c32)

	expected := map[int]string{
		0: "A Child 3 Level 2",
		1: "C Child 1 Level 2",
		2: "B Child 2 Level 1",
	}

	roles := Roles{&c12, &c21, &c32}
	sort.Sort(roles)

	for k, v := range roles {
		if expected[k] != v.Name {
			t.Errorf("`%s` is at position %d, should be `%s`", v.Name, k, expected[k])

		}
	}

}

func TestRoleCreateCircle(t *testing.T) {
	p := Role{Name: "Parent"}
	c1 := Role{Name: "Child Level 1"}
	c2 := Role{Name: "Child Level 2"}

	p.LinkChild(&c1)
	c1.LinkChild(&c2)

	err := c2.LinkChild(&p)
	if err == nil {
		t.Error("circular reflecton created")
	}
}

func TestRoleMultipleParents(t *testing.T) {
	p1 := Role{Name: "Parent 1"}
	p2 := Role{Name: "Parent 2"}
	c := Role{Name: "Child 1"}

	p1.LinkChild(&c)

	err := p2.LinkChild(&c)
	if err == nil {
		t.Error("child with parent added to new parent")
	}
}

func TestRoleDeleteChildFromParent(t *testing.T) {
	p := Role{Name: "Parent"}
	c := Role{Name: "Child"}

	p.LinkChild(&c)

	err := p.UnlinkChild(&c)
	if err != nil {
		t.Error("child not deleted from parent")
	}

	for _, elem := range p.Children {
		if elem == &c {
			t.Error("child not removed from parent")
		}
	}

	if c.Parent == &p {
		t.Error("child stil refers to parent")

	}
}

func TestRoleDeleteUnrelatedChild(t *testing.T) {
	p := Role{Name: "Parent"}
	c1 := Role{Name: "Child 1"}
	c2 := Role{Name: "Child 2"}

	p.LinkChild(&c1)
	err := p.UnlinkChild(&c2)
	if err == nil {
		t.Error("removing a node from another node that is not related as parent shuld be impossible")
	}
}
