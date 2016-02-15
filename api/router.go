package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/unprofession-al/gerty/entities"
)

var (
	ni entities.NodeInteractor
	ri entities.RoleInteractor
	rt http.Handler
)

func NewRouter(nodeInt entities.NodeInteractor, roleInt entities.RoleInteractor) http.Handler {
	ni = nodeInt
	ri = roleInt

	r := mux.NewRouter().StrictSlash(true)
	apiv1 := r.PathPrefix("/api/v1/").Subrouter()

	for name, route := range apiv1Routes {
		h := route.h
		if h == nil {
			h = notImplemented
		}

		apiv1.
			Methods(route.m).
			Path(route.p).
			Name(name).
			Handler(h)
	}

	chain := alice.New(recoverMiddleware, userContextMiddleware).Then(r)
	rt = chain

	return rt
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	user := getUserContext(r)
	w.WriteHeader(http.StatusNotImplemented)
	out := "Function Not Yet Implemented, " + user + "\n"
	w.Write([]byte(out))
}
