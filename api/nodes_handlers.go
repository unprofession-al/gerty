package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/helpers"
)

func listNodes(res http.ResponseWriter, req *http.Request) {
	out, err := ni.List()
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, out)
}

func getNode(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, node)
}

func addNode(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err == nil {
		helpers.Respond(res, req, http.StatusConflict, "already exists")
		return
	}

	node = entities.Node{Name: vars["node"]}

	err = ni.Save(node)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusCreated, node)
}

func delNode(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	err = ni.Delete(node)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusNoContent, "deleted")
}

func addNodeVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	var nodeVars map[string]interface{}
	err = helpers.ParseBodyAsMap(req, &nodeVars)
	if err != nil {
		helpers.Respond(res, req, http.StatusBadRequest, err.Error())
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
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusCreated, node)
}

func triggerNodeVarsProviders(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	buckets, err := np.GetVars(node.Name)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	for _, bucket := range buckets {
		node.Vars.AddOrReplaceBucket(bucket)
	}

	err = ni.Save(node)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusCreated, node)
}

func getNodeVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	out := ni.GetVars(node)

	helpers.Respond(res, req, http.StatusOK, out)
}

func getNodeVar(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	nvs := ni.GetVars(node)

	var out interface{}

	for _, nv := range nvs {
		if nv.Key == vars["var"] {
			out = nv.Value
			helpers.Respond(res, req, http.StatusOK, out)
			return
		}
	}

	helpers.Respond(res, req, http.StatusNotFound, out)
}

func replaceNodeVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	var nodeVars map[string]interface{}
	err = helpers.ParseBodyAsMap(req, &nodeVars)
	if err != nil {
		helpers.Respond(res, req, http.StatusBadRequest, err.Error())
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
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, node)
}

func linkNodeToRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node := entities.Node{Name: vars["node"]}
	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	for _, roleName := range node.Roles {
		if roleName == role.Name {
			helpers.Respond(res, req, http.StatusNotModified, node)
			return
		}
	}

	node.Roles = append(node.Roles, role.Name)

	err = ni.Save(node)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusCreated, node)
}

func unlinkNodeFromRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	node := entities.Node{Name: vars["node"]}
	node, err := ni.Get(vars["node"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
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
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, node)
}
