package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/helpers"
)

func listRoles(res http.ResponseWriter, req *http.Request) {
	out, err := ri.List()
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, out)
}

func getRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, role)
}

func addRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err == nil {
		helpers.Respond(res, req, http.StatusConflict, "already exists")
		return
	}

	role = entities.Role{Name: vars["role"]}

	err = ri.Save(role)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusCreated, role)
}

func delRole(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	// TODO: the relations must be sanitized on store level and must not be
	// done via handlers. Also, if a role is deleted, related nodes must be
	// unlinked.
	if role.Parent != "" {
		parent, err := ri.Get(role.Parent)
		if err != nil {
			helpers.Respond(res, req, http.StatusNotFound, err.Error())
			return
		}

		err = ri.UnlinkChild(&parent, &role)
		if err != nil {
			helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
			return
		}
	}

	for _, childName := range role.Children {
		child, err := ri.Get(childName)
		if err != nil {
			helpers.Respond(res, req, http.StatusNotFound, err.Error())
			return
		}

		err = ri.UnlinkChild(&role, &child)
		if err != nil {
			helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = ri.Delete(role)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.Respond(res, req, http.StatusOK, "deleted")
}

func addRoleVars(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	var roleVars map[string]interface{}
	err = helpers.ParseBodyAsMap(req, &roleVars)
	if err != nil {
		helpers.Respond(res, req, http.StatusBadRequest, err.Error())
		return
	}

	v := entities.VarBucket{
		Prio: 0,
		Vars: roleVars,
		Name: "native",
	}

	role.Vars.AddOrReplaceBucket(v)

	err = ri.Save(role)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusCreated, role)
}

func addRoleParent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	parent, err := ri.Get(vars["parent"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	err = ri.LinkParent(&role, &parent)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, role.Parent)
}

func delRoleParent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	role, err := ri.Get(vars["role"])
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	parent, err := ri.Get(role.Parent)
	if err != nil {
		helpers.Respond(res, req, http.StatusNotFound, err.Error())
		return
	}

	err = ri.UnlinkChild(&parent, &role)
	if err != nil {
		helpers.Respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Respond(res, req, http.StatusOK, "unlinked")
}
