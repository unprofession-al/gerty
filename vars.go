package main

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
		for k, v := range buck.Vars {
			merged.append(k, v, src, buck.Name)
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
	Key           string
	Value         interface{}
	Source        string
	SourceBucket  string
	Overwritten   bool
	OverwrittenBy *Merged
}

type MergedVars []*Merged

func (m MergedVars) String() string {
	var out string

	for _, v := range m {
		if v.Overwritten {
			out += "OVERWRITTEN BY `" + v.OverwrittenBy.Source + "/" + v.SourceBucket + "`: "
		}
		out += fmt.Sprintf("Key: `%s`, Value: {%s}, Source: `%s/%s` \n", v.Key, v.Value, v.Source, v.SourceBucket)
	}

	return out
}

func (m *MergedVars) append(key string, value interface{}, source string, bucket string) {
	v := &Merged{
		Key:          key,
		Value:        value,
		Source:       source,
		SourceBucket: bucket,
		Overwritten:  false,
	}

	for _, mv := range *m {
		if mv.Key == key && !mv.Overwritten {
			mv.OverwrittenBy = v
			mv.Overwritten = true
		}
	}

	*m = append(*m, v)
}
