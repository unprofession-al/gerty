package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sontags/env"
	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/entities"
	mw "github.com/unprofession-al/gerty/middleware"
	"github.com/unprofession-al/gerty/store"
	"github.com/unprofession-al/gerty/transformers"
	// _ "github.com/unprofession-al/gerty/store/memstore"
	_ "github.com/unprofession-al/gerty/store/sqlitestore"
)

type configuration struct {
	Port    string `json:"port"`
	Address string `json:"address"`
}

var config configuration

func init() {
	env.Var(&config.Port, "PORT", "8008", "Port to bind to")
	env.Var(&config.Address, "ADDR", "0.0.0.0", "Address to bind to")
}

func main() {
	env.Parse("GERTY", false)

	s, err := store.New("sqlitestore", "/tmp/gerty.sqlite3")
	if err != nil {
		panic(err)
	}

	ri := entities.NewRoleInteractor(s.Roles)
	ni := entities.NewNodeInteractor(s.Nodes, ri)

	r := mux.NewRouter().StrictSlash(true)

	api.Inject(ni, ri)
	a := r.PathPrefix("/api/").Subrouter()
	api.PopulateRouter(a)

	transformers.Inject(ni, ri)
	t := r.PathPrefix("/transformers/").Subrouter()
	transformers.PopulateRouter(t)

	chain := alice.New(
		mw.RecoverPanic,
		mw.CorsHeaders,
		mw.UserContext,
		// mw.WebHook,
	).Then(r)

	log.Fatal(http.ListenAndServe(config.Address+":"+config.Port, chain))
}
