package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
)

func listNodes(res http.ResponseWriter, req *http.Request) {
	out, err := ni.List()
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusOK, out)
}

func getNode(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	respond(res, req, http.StatusCreated, node)
}

func addNode(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var nodeVars map[string]interface{}
	err := parseBody(req, &nodeVars)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	v := entities.VarBucket{
		Prio: 0,
		Vars: nodeVars,
		Name: "native",
	}

	node := entities.Node{Name: vars["node"]}
	node.Vars.AddOrReplaceBucket(v)

	err = ni.Save(node)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusCreated, node)
}

func getNodeVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	out := ni.GetVars(node)

	respond(res, req, http.StatusOK, out)
}

func linkNodeToRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node := entities.Node{Name: vars["node"]}
	node, err := ni.Get(vars["node"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	role, err := ri.Get(vars["role"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	for _, roleName := range node.Roles {
		if roleName == role.Name {
			respond(res, req, http.StatusNotModified, node)
			return
		}
	}

	node.Roles = append(node.Roles, role.Name)

	err = ni.Save(node)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusCreated, node)
}
