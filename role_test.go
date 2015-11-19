package main

import "testing"

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
