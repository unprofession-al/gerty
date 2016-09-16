package providers

import (
	"net/http"
	"strings"

	"github.com/unprofession-al/gerty/entities"
)

type NodeVarsProvider struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Prio int    `json:"prio"`
}

func (np NodeVarsProvider) Get(nodename string) (entities.VarBucket, error) {
	var vars map[string]interface{}
	out := entities.VarBucket{
		Name: np.Name,
		Prio: np.Prio,
		Vars: vars,
	}

	url := strings.Replace(npr.Url, "{{nodename}}", nodename, -1)
	resp, err := http.Get(url)
	if err != nil {
		return out, err
	}

	err = getJSONBodyAsStruct(resp.Body, &np.Vars)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return out, err
	}

	return out, nil
}
