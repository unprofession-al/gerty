// Package entities provides all objects on which gerty operates on.
package entities

import "sort"

// Node representes a configuration item (eg. a server, network component etc.).
type Node struct {
	Name  string
	Vars  VarCollection
	Roles Roles
}

// GetVars collects and merges all defined variables that are relevant to a node.
// The variables returned consist of all definend variables of the node itself as
// well as the variables of its roles. The following steps are performed:
//
// 1) All roles to which the node is assigned are sorted: Deep to flat, roles with
// the same depth in alphabetical order.
//
// 2) Starting from the first roles variables are merged towards the root role.
// Variables found first win since they are closer to the node itself.
// Merging of the branch is stopped as soon as either the root element is reached
// or a role is reached that has already been collected before.
//
// 3) Step 2) is repeated for each branch. Already collected roles are tracked for
// all branches together.
//
// 4) The variables of each branch are merged. The variable with the smallest
// distance to the node will win. A variable is defined in multiple branches with
// the same distance to the node, merging is unambiguous. The variable of the role
// that comes first in the alphabet will win. Also the losing variable will be
// referenced in the `tainting` field of the winning variable.
//
// 5) Node variables will be merged with the consolideted role variables. Node
// always win.
func (n Node) GetVars() MergedVars {
	visited := []*Role{}
	sort.Sort(n.Roles)
	merged := MergedVars{}

	for _, role := range n.Roles {
		branchVars := MergedVars{}
		distance := 0
		current := role
		for true {
			if contains(visited, current) {
				break
			}
			distance++
			for _, v := range current.Vars.Merge(current.Name) {
				v.Distance = distance
				branchVars.insertAsOldest(v)
			}
			visited = append(visited, current)
			if current.GetParent() != nil {
				current = current.GetParent()
			} else {
				break
			}
		}

		for _, v := range branchVars {
			merged.insertNearestAsNewest(v)
		}
	}

	for _, v := range n.Vars.Merge("Node") {
		merged.insertAsNewest(v)
	}

	return merged

}

func contains(roles []*Role, r *Role) bool {
	for _, role := range roles {
		if role == r {
			return true
		}
	}
	return false
}