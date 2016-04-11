package transformers

import (
	"net/http"

	"github.com/unprofession-al/gerty/helpers"
)

type AnsibleInventory map[string]*AnsibleGroup

type AnsibleGroup struct {
	Hosts    []string                          `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Vars     map[string]interface{}            `json:"vars,omitempty" yaml:"vars,omitempty"`
	Hostvars map[string]map[string]interface{} `json:"hostvars,omitempty" yaml:"hostvars,omitempty"`
	Children []string                          `json:"childern,omitempty" yaml:"childern,omitempty"`
}

func NewAnsibleInventory() (*AnsibleInventory, error) {
	i := AnsibleInventory{}

	metaGroupName := "_meta"
	i[metaGroupName] = &AnsibleGroup{}
	i[metaGroupName].Hostvars = make(map[string]map[string]interface{})

	nodes, err := ni.List()
	if err != nil {
		return &i, err
	}

	roles, err := ri.List()
	if err != nil {
		return &i, err
	}

	for _, rName := range roles {

		r, err := ri.Get(rName)
		if err != nil {
			return &i, err
		}
		ag := &AnsibleGroup{}
		ag.Children = r.Children
		i[r.Name] = ag
	}

	for _, nName := range nodes {
		n, err := ni.Get(nName)
		if err != nil {
			return &i, err
		}

		for _, r := range n.Roles {
			i[r].Hosts = append(i[r].Hosts, n.Name)
		}

		merged := ni.GetVars(n)
		hostvars := make(map[string]interface{})
		for _, v := range merged {
			hostvars[v.Key] = v.Value
		}
		i[metaGroupName].Hostvars[n.Name] = hostvars
	}
	return &i, nil
}

func ansibleRenderer(res http.ResponseWriter, req *http.Request) {
	out, err := NewAnsibleInventory()
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.Respond(res, req, http.StatusOK, out)
}
