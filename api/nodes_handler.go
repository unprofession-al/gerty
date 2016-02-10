package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unrolled/render"
)

func listNodes(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	out := ni.List()
	r.JSON(res, http.StatusOK, out)
}

func getNode(res http.ResponseWriter, req *http.Request) {
	r := render.New()

	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	r.JSON(res, http.StatusOK, node)
}

func addNode(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	node := entities.Node{Name: vars["node"]}

	err := ni.Save(node)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	r.JSON(res, http.StatusCreated, node)
}

func getNodeVars(res http.ResponseWriter, req *http.Request) {
	r := render.New()

	vars := mux.Vars(req)

	node, err := ni.Get(vars["node"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	out := ni.GetVars(node)

	r.JSON(res, http.StatusOK, out)
}

func linkNodeToRole(res http.ResponseWriter, req *http.Request) {
	r := render.New()

	vars := mux.Vars(req)

	node := entities.Node{Name: vars["node"]}
	node, err := ni.Get(vars["node"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	role, err := ri.Get(vars["role"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	for _, roleName := range node.Roles {
		if roleName == role.Name {
			r.JSON(res, http.StatusNotModified, node)
			return
		}
	}

	node.Roles = append(node.Roles, role.Name)

	err = ni.Save(node)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	r.JSON(res, http.StatusCreated, node)
}
