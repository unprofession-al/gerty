package entities

import (
	"fmt"
	"testing"
)

var tests = []VarBucket{
	{
		Name: "Bucket 1",
		Prio: 1,
		Vars: VarList{
			"Key 1": "Value A",
			"Key 2": "Value B",
		},
	},
	{
		Name: "Bucket 2",
		Prio: 0,
		Vars: VarList{
			"Key 2": "Value C",
			"Key 3": "Value D",
		},
	},
	{
		Name: "Bucket 3",
		Prio: 3,
		Vars: VarList{
			"Key 2": "Value E",
			"Key 3": "Value F",
			"Key 4": "Value G",
		},
	},
}

var results = map[string]interface{}{
	"Key 1": "Value A",
	"Key 2": "Value C",
	"Key 3": "Value D",
	"Key 4": "Value G",
}

func TestVarsSerialize(t *testing.T) {
	c := VarCollection{}
	for _, bucket := range tests {
		c.AddOrReplaceBucket(bucket)
	}
	b, err := c.Serialize()
	if err != nil {
		t.Error("Could not serialize")
	}
	fmt.Println(string(b))
}

func TestVarsMerge(t *testing.T) {
	c := VarCollection{}
	for _, bucket := range tests {
		c.AddOrReplaceBucket(bucket)
	}

	merged := c.Merge("Test")
	fmt.Println(merged)
	for rk, rv := range results {
		found := false
		for _, v := range merged {
			if v.Key == rk {
				found = true
				if v.Value != rv {
					t.Errorf("Var `%s` has value `%s`, should have `%s`", rk, v.Value, results[v.Key])
				}
			}
		}
		if !found {
			t.Errorf("Var `%s` should exist but does not", rk)
		}
	}
}

func TestVarsAddReplace(t *testing.T) {
	c := VarCollection{}
	for _, bucket := range tests {
		c.AddOrReplaceBucket(bucket)
		c.AddOrReplaceBucket(bucket) // add bucket twice
	}

	expected := len(tests)

	if len(c) < expected {
		t.Error("VarCollection contains more elements than expected")
	} else if len(c) > expected {
		t.Error("VarCollection contains less elements than expected")
	}
}
