package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unrolled/render"
)

func listRoles(res http.ResponseWriter, req *http.Request) {
	r := render.New()

	out, err := ri.List()
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	r.JSON(res, http.StatusOK, out)
}

func getRole(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	r.JSON(res, http.StatusOK, role)
}

func addRole(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	role := entities.Role{Name: vars["role"]}

	err := ri.Save(role)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	r.JSON(res, http.StatusCreated, role)
}

func getRoleParent(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	r.JSON(res, http.StatusOK, role.Parent)
}

func addRoleParent(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	parent, err := ri.Get(vars["parent"])
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	err = ri.LinkChild(&parent, &role)
	if err != nil {
		r.JSON(res, http.StatusNotFound, err.Error())
		return
	}

	r.JSON(res, http.StatusOK, role.Parent)
}
