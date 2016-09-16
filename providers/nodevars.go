package providers

// use as
//	providersConfig := `[{"prio":3, "name": "netmgmt", "url": "http://netmgmt.stxt.media.int:9001/nodes/{{nodename}}"}]`
//	p, _ := providers.Bootstrap(providersConfig)
//	if err != nil {
//		panic(err)
//	}
//	nodevars, err := p.GetVars("devops-01.stxt.media.int")
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(nodevars)

import (
	"encoding/json"
	"fmt"
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
	Name string `json:"name"`
	Url  string `json:"url"`
	Prio int    `json:"prio"`
}

func (np NodeVarsProvider) GetVars(nodename string) (entities.VarBucket, error) {
	out := entities.VarBucket{
		Name: np.Name,
		Prio: np.Prio,
	}

	url := strings.Replace(np.Url, "{{nodename}}", nodename, -1)
	fmt.Println(url)
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
	out.Vars = vars

	return out, nil
}
