package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
)

func listRoles(res http.ResponseWriter, req *http.Request) {
	out, err := ri.List()
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	respond(res, req, http.StatusOK, out)
}

func getRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	respond(res, req, http.StatusOK, role)
}

func addRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role := entities.Role{Name: vars["role"]}

	err := ri.Save(role)
	if err != nil {
		respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, req, http.StatusCreated, role)
}

func getRoleParent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	respond(res, req, http.StatusOK, role.Parent)
}

func addRoleParent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	parent, err := ri.Get(vars["parent"])
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	err = ri.LinkChild(&parent, &role)
	if err != nil {
		respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	respond(res, req, http.StatusOK, role.Parent)
}
