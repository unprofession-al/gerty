package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/entities"
	mw "github.com/unprofession-al/gerty/middleware"
	"github.com/unprofession-al/gerty/store"
	_ "github.com/unprofession-al/gerty/store/memstore"
	_ "github.com/unprofession-al/gerty/store/sqlitestore"
)

func main() {
	s, err := store.New("sqlitestore", "/tmp/gerty.sqlite3")
	if err != nil {
		panic(err)
	}

	ri := entities.NewRoleInteractor(s.Roles)
	ni := entities.NewNodeInteractor(s.Nodes, ri)

	api.InjectRouter(ni, ri)
	r := mux.NewRouter().StrictSlash(true)
	api.PopulateRouter(r)

	chain := alice.New(mw.RecoverPanic, mw.UserContext).Then(r)
	//apiRouter := api.NewRouter(ni, ri)

	log.Fatal(http.ListenAndServe(":8008", chain))
}
