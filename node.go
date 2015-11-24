package main

import "sort"

type Node struct {
	ID    int
	Name  string
	Vars  VarCollection
	Roles Roles
}

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
				branch_vars.InsertAsOldest(v.Key, v.Value, v.Source, v.SourceBucket, distance)
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
		merged.InsertAsNewest(v.Key, v.Value, v.Source, v.SourceBucket, 0)
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
