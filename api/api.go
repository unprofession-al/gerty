package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/gerty/config"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/providers"
)

var (
	ni  entities.NodeInteractor
	ri  entities.RoleInteractor
	np  providers.NodeVarsProviders
	cfg config.Configuration
)

// Inject stores the required Interactors as global vars in order to
// provide access to the persistence layer and business logic.
func Inject(nodeInt entities.NodeInteractor, roleInt entities.RoleInteractor) {
	ni = nodeInt
	ri = roleInt
}

func Configure(config config.Configuration) error {
	cfg = config

	var err error
	np, err = providers.Bootstrap(cfg.NodeVarsProviders)
	return err
}

func notImplemented(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotImplemented)
	out := "Function Not Yet Implemented\n"
	res.Write([]byte(out))
}

// PopulateRouter appends all defined routes to a given gorilla mux router.
func PopulateRouter(router *mux.Router) {
	for version, leafs := range routes {
		api := router.PathPrefix("/" + version).Subrouter()

		for pattern, leaf := range leafs {
			appendLeaf(pattern, leaf, api)
		}
	}
}

func appendLeaf(p string, l leaf, router *mux.Router) {
	for method, endpoint := range l.E {
		h := endpoint.h
		if h == nil {
			h = notImplemented
		}
		router.
			Methods(method).
			Path("/" + p).
			Name(endpoint.N).
			Handler(h)
	}
	for pattern, leaf := range l.L {
		appendLeaf(p+"/"+pattern, leaf, router)
	}
}

var routes = make(map[string]leafs)

type leafs map[string]leaf

type leaf struct {
	E endpoints `json:"endpoints,omitempty"`
	L leafs     `json:"leafs,omitempty"`
}

type endpoints map[string]endpoint

type endpoint struct {
	N string `json:"name"`
	h http.HandlerFunc
}
