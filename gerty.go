package main

import (
	"log"
	"net/http"

	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/entities"
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

	apiRouter := api.NewRouter(ni, ri)

	log.Fatal(http.ListenAndServe(":8008", apiRouter))
}
