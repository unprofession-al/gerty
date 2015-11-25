// Package entities provides all objects on which gerty operates on.
package entities

import "sort"

//
type Node struct {
	ID    int
	Name  string
	Vars  VarCollection
	Roles Roles
}

// GetVars collects and merges all defined variables that are relevant to a node.
// In order to have a deterministic way to do so, the following steps are performed:
//
// - The Roles a node is assigned to are sorted
//
func (n Node) GetVars() MergedVars {
	visited := []int{}
	sort.Sort(n.Roles)
	merged := MergedVars{}

	for _, role := range n.Roles {
		branch_vars := MergedVars{}
		distance := 0
		current := role
		for true {
			if contains(visited, current.ID) {
				break
			}
			distance += 1
			for _, v := range current.Vars.Merge(current.Name) {
				v.Distance = distance
				branch_vars.InsertAsOldest(v)
			}
			visited = append(visited, current.ID)
			if current.Parent != nil {
				current = current.Parent
			} else {
				break
			}
		}

		for _, v := range branch_vars {
			merged.InsertNearer(v)
		}
	}

	for _, v := range n.Vars.Merge("Node") {
		merged.InsertAsNewest(v)
	}

	return merged

}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
