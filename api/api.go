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

var routes = make(map[string]routeDefinition)

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

// PopulateRouter appends all defined routes to a given gorilla mux router.
func PopulateRouter(router *mux.Router) {
	for version, r := range routes {
		api := router.PathPrefix("/" + version + "/").Subrouter()

		for name, route := range r {
			h := route.h
			if h == nil {
				h = notImplemented
			}

			api.
				Methods(route.M).
				Path(route.P).
				Name(name).
				Handler(h)
		}
	}
}

func notImplemented(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotImplemented)
	out := "Function Not Yet Implemented\n"
	res.Write([]byte(out))
}

type route struct {
	// method
	M string `json:"method"`
	// pattern
	P string `json:"pattern"`
	// HandlerFunc
	h http.HandlerFunc
}

type routeDefinition map[string]route
