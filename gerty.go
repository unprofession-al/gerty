package main

import (
	"log"
	"net/http"

	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/store"
)

func main() {
	stores, err := store.Open("mem", "")
	if err != nil {
		panic(err)
	}

	ri := entities.NewRoleInteractor(stores.Roles)
	ni := entities.NewNodeInteractor(stores.Nodes, ri)

	router := api.NewRouter(ni, ri)

	log.Fatal(http.ListenAndServe(":8008", router))
}
