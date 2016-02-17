package entities

import (
	"encoding/json"
	"fmt"
	"sort"
)

// VarBucket adds a name and a priority to a VarList.
type VarBucket struct {
	Name string                 `json:"name"`
	Prio int                    `json:"prio"` // lower value means higher priority
	Vars map[string]interface{} `json:"vars"`
}

// VarCollection groups together multiple VarBuckets. This entity is
// referenced by nodes and roles.
type VarCollection []VarBucket

// Len returns the length of the list, part of implementing the sort interface.
func (v VarCollection) Len() int { return len(v) }

// Swap changes the order of two elements in the list, part of implementing
// the sort interface.
func (v VarCollection) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// Less compares the order of two elements, part of implementing the sort interface.
// The lower the value of `Prio` the higher the priority.
func (v VarCollection) Less(i, j int) bool { return v[i].Prio > v[j].Prio }

// Merge consolidates all variables devined in the VarBuckets. Buckets with the
// higher priority win.
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
			merged.insertAsNewest(v)
		}
	}
	return merged
}

// AddOrReplaceBucket adds a VarBucket to a collection. If a bucket with the same
// name exists, the existing Bucket will be replaced.
func (v *VarCollection) AddOrReplaceBucket(b VarBucket) {
	for i, bucket := range *v {
		if bucket.Name == b.Name {
			(*v)[i] = b
			return
		}
	}
	*v = append(*v, b)
	return
}

// Serialize renders the colletion as json byte slice
func (v *VarCollection) Serialize() ([]byte, error) {
	return json.Marshal(v)
}

// Deserialize conberts a byte slice to a VarCollection
func (v *VarCollection) Deserialize(b []byte) error {
	return json.Unmarshal(b, v)
}

// Merged holds a key/value representation of a variable as well as some meta data
// about the variable.
type Merged struct {
	// key of the variable
	Key string
	// value of the variable
	Value interface{}
	// name of the source where name usually represents the name of the role/node containing the VarCollection
	Source string
	// bucket name inside the containing VarCollection
	SourceBucket string
	// reference to Merged representations of earlier states of the same variable
	Old *Merged
	// distance metween the variable origin and the merging element (generally an node)
	Distance int
	// if merging is unambiguous, the variable that tries to merge later will be referenced as `tainting`
	Tainting *Merged
}

// String implements the `Stringer` interface for debugging reasons
func (m Merged) String() string {
	tainting := ""
	if m.Tainting != nil {
		tainting = fmt.Sprintf(", Tainting: `%s/%s`", m.Tainting.Source, m.Tainting.SourceBucket)
	}
	return fmt.Sprintf("Key: `%s`, Value: {%s}, Source: `%s/%s`, Dist: %d%s \n", m.Key, m.Value, m.Source, m.SourceBucket, m.Distance, tainting)
}

// MergedVars is a list of references to `Merged` elements
type MergedVars []*Merged

// String implements the `Stringer` interface for debugging reasons
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

func (m *MergedVars) insertAsNewest(v *Merged) {
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

func (m *MergedVars) insertNearestAsNewest(v *Merged) {
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

func (m *MergedVars) insertAsOldest(v *Merged) {
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
