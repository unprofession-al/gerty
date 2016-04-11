package entities

import (
	"fmt"
	"testing"
)

func TestNodeInteractor(t *testing.T) {
	newNode := Node{Name: "TestNode"}

	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})
	ni := NewNodeInteractor(NodeStoreMock{nodes: make(map[string]Node)}, ri)

	ni.Save(newNode)

	node, _ := ni.Get(newNode.Name)

	if node.Name != newNode.Name {
		t.Errorf("Name in not consistent: `%s` != '%s'", node.Name, newNode.Name)
	}
}

func TestListNodes(t *testing.T) {
	node1 := Node{Name: "Node 1"}
	node2 := Node{Name: "Node 2"}
	node3 := Node{Name: "Node 3"}

	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})
	ni := NewNodeInteractor(NodeStoreMock{nodes: make(map[string]Node)}, ri)

	ni.Save(node1)
	ni.Save(node2)
	ni.Save(node3)

	nodeList, _ := ni.List()

	if len(nodeList) != 3 {
		t.Errorf("Wrong number of nodes found: `%d` != '%d'", len(nodeList), 3)
	}
}

/*
	           ra
                |
        +-------+-------+
		|               |
	   rb              rc
        |               |
    +---+---+       +---+
    |       |       |
   rd       re      rf
    |
+---+---+
|       |
rg     rh
*/

var r = map[string]*Role{
	"a": {
		Name: "ra",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 1": "Value A",
					"Var 2": "Value A",
					"Var 3": "Value A",
					"Var 4": "Value A",
					"Var 5": "Value A",
					"Var 7": "Value A",
					"Var 8": "Value A",
				},
			},
		},
	},
	"b": {
		Name: "rb",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 2": "Value B",
				},
			},
		},
	},
	"c": {
		Name: "rc",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 4": "Value C",
					"Var 5": "Value C",
				},
			},
		},
	},
	"d": {
		Name: "rd",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 1": "Value D",
				},
			},
		},
	},
	"e": {
		Name: "re",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 3": "Value E",
					"Var 6": "Value E",
				},
			},
		},
	},
	"f": {
		Name: "rf",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 5": "Value F",
					"Var 6": "Value F",
				},
			},
		},
	},
	"g": {
		Name: "rg",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 3": "Value G",
				},
			},
		},
	},
	"h": {
		Name: "rh",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 1": "Value H",
				},
			},
		},
	},
}

var nodeResults = map[string]interface{}{
	"Var 1": "Value H",
	"Var 2": "Value B",
	"Var 3": "Value G",
	"Var 4": "Value C",
	"Var 5": "Value F",
	"Var 6": "Value E",
	"Var 7": "Value A",
	"Var 8": "Value HOST",
}

var tainting = map[string]bool{
	"Var 1": false,
	"Var 2": false,
	"Var 3": true,
	"Var 4": false,
	"Var 5": false,
	"Var 6": true,
	"Var 7": false,
	"Var 8": false,
}

func TestNodeMerging(t *testing.T) {
	ri := NewRoleInteractor(RoleStoreMock{roles: make(map[string]Role)})
	for _, role := range r {
		ri.Save(*role)
	}

	ri.LinkParent(r["b"], r["a"])
	ri.LinkParent(r["c"], r["a"])
	ri.LinkParent(r["d"], r["b"])
	ri.LinkParent(r["e"], r["b"])
	ri.LinkParent(r["f"], r["c"])
	ri.LinkParent(r["g"], r["d"])
	ri.LinkParent(r["h"], r["d"])

	ni := NewNodeInteractor(NodeStoreMock{nodes: make(map[string]Node)}, ri)

	roles := []string{r["g"].Name, r["h"].Name, r["e"].Name, r["f"].Name, r["c"].Name}
	node := Node{
		Name:  "Test",
		Roles: roles,
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: map[string]interface{}{
					"Var 8": "Value HOST",
				},
			},
		},
	}

	ni.Save(node)

	vars := ni.GetVars(node)
	fmt.Println(vars)
	for rk, rv := range nodeResults {
		found := false
		for _, v := range vars {
			if v.Key == rk {
				found = true
				if v.Value != rv {
					t.Errorf("Var `%s` has value `%s`, should have `%s`", rk, v.Value, nodeResults[v.Key])
				}
				taint := false
				if v.Tainting != nil {
					taint = true
				}
				if taint != tainting[rk] {
					t.Errorf("Var `%s` has tainted `%v`, should have `%v`", rk, taint, tainting[rk])
				}
			}
		}
		if !found {
			t.Errorf("Var `%s` should exist but does not", rk)
		}
	}
}
