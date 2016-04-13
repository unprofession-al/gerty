package transformers

import (
	"net/http"

	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/helpers"
)

type treeNode struct {
	Name     string      `json:"name"`
	Children []*treeNode `json:"children"`
	Type     string      `json:"type"`
}

func treeRenderer(res http.ResponseWriter, req *http.Request) {
	rootRole := ""
	root := req.URL.Query()["role"]
	if len(root) > 0 {
		rootRole = root[0]
	}

	out := &treeNode{}

	role, err := ri.Get(rootRole)
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	addToTree(role, out)

	helpers.Respond(res, req, http.StatusCreated, out)
}

func addToTree(current entities.Role, t *treeNode) {
	t.Name = current.Name
	t.Type = "role"
	for _, host := range current.Hosts {
		h := &treeNode{
			Name: host,
			Type: "host",
		}
		t.Children = append(t.Children, h)
	}
	for _, child := range current.Children {
		role, _ := ri.Get(child)

		childTree := &treeNode{}
		addToTree(role, childTree)
		t.Children = append(t.Children, childTree)
	}
}
