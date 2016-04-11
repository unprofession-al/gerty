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

	node, err := ni.Get(vars["node"])
	if err == nil {
		respond(res, req, http.StatusConflict, "already exists")
		return
	}

	node = entities.Node{Name: vars["node"]}

	err = ni.Save(node)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusCreated, node)
}

func delNode(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	err = ni.Delete(node)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusOK, "deleted")
}

func addNodeVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	var nodeVars map[string]interface{}
	err = parseBodyAsMap(req, &nodeVars)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	v := entities.VarBucket{
		Prio: 0,
		Vars: nodeVars,
		Name: "native",
	}

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

func replaceNodeVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	var nodeVars map[string]interface{}
	err = parseBodyAsMap(req, &nodeVars)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	vb := entities.VarBucket{
		Prio: 0,
		Name: "native",
	}
	for _, b := range node.Vars {
		if b.Name == "native" {
			vb = b
		}
	}

	for key, val := range nodeVars {
		vb.Vars[key] = val
	}

	node.Vars.AddOrReplaceBucket(vb)

	err = ni.Save(node)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusOK, node)
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

func unlinkNodeFromRole(res http.ResponseWriter, req *http.Request) {
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

	for i, roleName := range node.Roles {
		if roleName == role.Name {
			node.Roles = append(node.Roles[:i], node.Roles[i+1:]...)
			break
		}
	}

	err = ni.Save(node)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusOK, node)
}
