package main

import (
	"sort"
	"testing"
)

func TestRoleAddChildToParent(t *testing.T) {
	p := Role{ID: 1, Name: "Parent"}
	c := Role{ID: 2, Name: "Child"}

	err := p.AddChild(&c)
	if err != nil {
		t.Error("child not added to parent")
	}

	childrenAdded := false
	for _, elem := range p.Children {
		if elem.ID == c.ID {
			childrenAdded = true
		}
	}
	if !childrenAdded {
		t.Error("child not added to parent")
	}

	if c.Parent.ID != p.ID {
		t.Error("group not refered in child")

	}
}

func TestRoleCalculateDepth(t *testing.T) {
	p := Role{ID: 1, Name: "Parent"}
	c1 := Role{ID: 2, Name: "Child Level 1"}
	c2 := Role{ID: 3, Name: "Child Level 2"}

	err := p.AddChild(&c1)
	if err != nil {
		t.Error("child not added to parent")
	}

	err = c1.AddChild(&c2)
	if err != nil {
		t.Error("child not added to parent")
	}

	if c2.Depth() != 2 {
		t.Error("depth not calculated correctly")
	}
}

func TestRoleSort(t *testing.T) {
	p1 := Role{ID: 1, Name: "C Parent 1"}
	c11 := Role{ID: 2, Name: "C Child 1 Level 1"}
	c12 := Role{ID: 3, Name: "C Child 1 Level 2"}
	p1.AddChild(&c11)
	c11.AddChild(&c12)

	p2 := Role{ID: 4, Name: "B Parent 2"}
	c21 := Role{ID: 5, Name: "B Child 2 Level 1"}
	p2.AddChild(&c21)

	p3 := Role{ID: 6, Name: "A Parent 3"}
	c31 := Role{ID: 7, Name: "A Child 3 Level 1"}
	c32 := Role{ID: 8, Name: "A Child 3 Level 2"}
	p3.AddChild(&c31)
	c31.AddChild(&c32)

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
	p := Role{ID: 1, Name: "Parent"}
	c1 := Role{ID: 2, Name: "Child Level 1"}
	c2 := Role{ID: 3, Name: "Child Level 2"}

	err := p.AddChild(&c1)
	if err != nil {
		t.Error("child not added to parent")
	}

	err = c1.AddChild(&c2)
	if err != nil {
		t.Error("child not added to parent")
	}

	err = c2.AddChild(&p)
	if err == nil {
		t.Error("circular reflecton created")
	}
}

func TestRoleMultipleParents(t *testing.T) {
	p1 := Role{ID: 1, Name: "Parent 1"}
	p2 := Role{ID: 2, Name: "Parent 2"}
	c := Role{ID: 3, Name: "Child 1"}

	err := p1.AddChild(&c)
	if err != nil {
		t.Error("child not added to parent")
	}

	err = p2.AddChild(&c)
	if err == nil {
		t.Error("child with parent added to new parent")
	}
}

func TestRoleDeleteChildFromParent(t *testing.T) {
	p := Role{ID: 1, Name: "Parent"}
	c := Role{ID: 2, Name: "Child"}

	err := p.AddChild(&c)
	if err != nil {
		t.Error("child not added to parent")
	}

	err = p.DeleteChild(&c)
	if err != nil {
		t.Error("child not deleted from parent")
	}

	for _, elem := range p.Children {
		if elem.ID == c.ID {
			t.Error("child not removed from parent")
		}
	}

	if c.Parent.ID == p.ID {
		t.Error("child stil refers to parent")

	}
}

func TestRoleDeleteUnrelatedChild(t *testing.T) {
	p := Role{ID: 1, Name: "Parent"}
	c1 := Role{ID: 2, Name: "Child 1"}
	c2 := Role{ID: 3, Name: "Child 2"}

	err := p.AddChild(&c1)
	if err != nil {
		t.Error("child not added to parent")
	}

	err = p.DeleteChild(&c2)
	if err == nil {
		t.Error("removing a node from another node that is not related as parent shuld be impossible")
	}
}
