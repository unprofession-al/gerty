package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/entities"
)

var (
	ni entities.NodeInteractor
	ri entities.RoleInteractor
	rt *mux.Router
)

func NewRouter(nodeInt entities.NodeInteractor, roleInt entities.RoleInteractor) *mux.Router {
	ni = nodeInt
	ri = roleInt

	r := mux.NewRouter().StrictSlash(true)
	apiv1 := r.PathPrefix("/api/v1/").Subrouter()

	for name, route := range apiv1Routes {
		var handler http.Handler
		handler = route.HandlerFunc

		apiv1.
			Methods(route.Method).
			Path(route.Pattern).
			Name(name).
			Handler(handler)
	}

	rt = r

	return r
}

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes map[string]Route

var apiv1Routes = Routes{
	"ListNodes": Route{
		Method:      "GET",
		Pattern:     "/nodes/",
		HandlerFunc: listNodes,
	},
	"GetNode": Route{
		Method:      "GET",
		Pattern:     "/nodes/{node}",
		HandlerFunc: getNode,
	},
	"AddNode": Route{
		Method:      "POST",
		Pattern:     "/nodes/{node}",
		HandlerFunc: addNode,
	},
	"GetNodeVars": Route{
		Method:      "GET",
		Pattern:     "/nodes/{node}/vars",
		HandlerFunc: getNodeVars,
	},
	"LinkNodeToRole": Route{
		Method:      "POST",
		Pattern:     "/nodes/{node}/roles/{role}",
		HandlerFunc: linkNodeToRole,
	},
	"ListRoles": Route{
		Method:      "GET",
		Pattern:     "/roles/",
		HandlerFunc: listRoles,
	},
	"GetRole": Route{
		Method:      "GET",
		Pattern:     "/roles/{role}",
		HandlerFunc: getRole,
	},
	"AddRole": Route{
		Method:      "POST",
		Pattern:     "/roles/{role}",
		HandlerFunc: addRole,
	},
}
