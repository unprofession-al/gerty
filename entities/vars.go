package entities

import (
	"fmt"
	"sort"
)

type VarList map[string]interface{}

type VarBucket struct {
	Name string
	Prio int
	Vars VarList
}

type VarCollection []VarBucket

func (v VarCollection) Len() int           { return len(v) }
func (v VarCollection) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v VarCollection) Less(i, j int) bool { return v[i].Prio > v[j].Prio }

func (v *VarCollection) Merge(src string) MergedVars {
	sort.Sort(v)
	merged := MergedVars{}
	for _, buck := range *v {
		for k, val := range buck.Vars {
			v := &Merged{
				Key:          k,
				Value:        val,
				Source:       src,
				SourceBucket: buck.Name,
				Distance:     0,
			}
			merged.InsertAsNewest(v)
		}
	}
	return merged
}

func (v *VarCollection) AddOrReplace(n VarBucket) {
	for i, vars := range *v {
		if vars.Name == n.Name {
			(*v)[i] = n
			return
		}
	}
	*v = append(*v, n)
	return
}

type Merged struct {
	Key          string
	Value        interface{}
	Source       string
	SourceBucket string
	Old          *Merged
	Distance     int
	Tainting     *Merged
}

func (m Merged) String() string {
	tainting := ""
	if m.Tainting != nil {
		tainting = fmt.Sprintf(", Tainting: `%s/%s`", m.Tainting.Source, m.Tainting.SourceBucket)
	}
	return fmt.Sprintf("Key: `%s`, Value: {%s}, Source: `%s/%s`, Dist: %d%s \n", m.Key, m.Value, m.Source, m.SourceBucket, m.Distance, tainting)
}

type MergedVars []*Merged

func (m MergedVars) String() string {
	var out string

	for _, v := range m {
		history := ""
		current := v
		indent := ""
		for true {
			if current.Old != nil {
				indent += "-"
				history += indent + " " + current.Old.String()
				current = current.Old
			} else {
				break
			}
		}
		out += fmt.Sprintf("%s%s", v.String(), history)
	}

	return out
}

func (m *MergedVars) InsertAsNewest(v *Merged) {
	found := false
	for k, mv := range *m {
		if mv.Key == v.Key {
			v.Old = mv
			(*m)[k] = v
			found = true
			break
		}
	}

	if !found {
		*m = append(*m, v)
	}
}

func (m *MergedVars) InsertNearer(v *Merged) {
	found := false
	for k, mv := range *m {
		if mv.Key == v.Key {
			found = true
			if mv.Distance > v.Distance {
				v.Old = mv
				(*m)[k] = v
			} else if mv.Distance == v.Distance {
				mv.Tainting = v
			}
			break
		}
	}

	if !found {
		*m = append(*m, v)
	}
}

func (m *MergedVars) InsertAsOldest(v *Merged) {
	inserted := false
	for _, mv := range *m {
		if mv.Key == v.Key {
			oldest := mv
			for true {
				if oldest.Old != nil {
					oldest = oldest.Old
				} else {
					mv.Old = v
					inserted = true
					break
				}
			}
			break
		}
	}
	if !inserted {
		*m = append(*m, v)
	}
}
