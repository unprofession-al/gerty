package api

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	nodes := api.PathPrefix("/nodes").Subrouter()

	nodes.HandleFunc("/", listNodes).Methods("GET")
	return r
}
