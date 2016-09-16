package providers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/unprofession-al/gerty/entities"
)

type NodeVarsProviders []NodeVarsProvider

func Bootstrap(config string) (NodeVarsProviders, error) {
	nps := NodeVarsProviders{}
	err := json.Unmarshal([]byte(config), &nps)
	return nps, err
}

func (nps NodeVarsProviders) GetVars(nodename string) ([]entities.VarBucket, error) {
	buckets := []entities.VarBucket{}
	for _, np := range nps {
		bucket, err := np.GetVars(nodename)
		if err != nil {
			return buckets, err
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}

type NodeVarsProvider struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Prio   int    `json:"prio"`
	Prefix string `json:"prefix"`
}

func (np NodeVarsProvider) GetVars(nodename string) (entities.VarBucket, error) {
	out := entities.VarBucket{
		Name: np.Name,
		Prio: np.Prio,
	}

	url := strings.Replace(np.Url, "{{nodename}}", nodename, -1)
	resp, err := http.Get(url)
	if err != nil {
		return out, err
	}

	serialized, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	var vars map[string]interface{}
	err = json.Unmarshal(serialized, &vars)
	if err != nil {
		return out, err
	}

	if np.Prefix != "" {
		vars = prefix(vars, np.Prefix)
	}

	out.Vars = vars

	return out, nil
}

func prefix(vars map[string]interface{}, prefix string) map[string]interface{} {
	out := make(map[string]interface{})
	for name, val := range vars {
		out[prefix+"_"+name] = val
	}
	return out
}
