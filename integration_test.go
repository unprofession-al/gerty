package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/entities"
	"github.com/unprofession-al/gerty/store"
	_ "github.com/unprofession-al/gerty/store/memstore"
)

var (
	server *httptest.Server
	router http.Handler
	reader io.Reader
	nodes  []string
	roles  []string
)

func bootstrap() {
	nodes = []string{"node1", "node_2", "node-3"}
	roles = []string{"role1", "role_2", "role-3"}

	stores, err := store.New("memstore", "")
	if err != nil {
		panic(err)
	}

	ri := entities.NewRoleInteractor(stores.Roles)
	ni := entities.NewNodeInteractor(stores.Nodes, ri)

	router = api.NewRouter(ni, ri)
	server = httptest.NewServer(router)
}

func TestAddNodes(t *testing.T) {
	bootstrap()
	for _, node := range nodes {
		path, err := router.Get("AddNode").URL("node", node)
		if err != nil {
			t.Error(err)
		}

		url := server.URL + path.String()
		res, err := http.Post(url, "application/json", nil)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != http.StatusCreated {
			t.Errorf("Success expected at %s: %d", url, res.StatusCode)
		}
	}
}

func TestAddRoles(t *testing.T) {
	bootstrap()
	for _, role := range roles {
		path, err := router.Get("AddRole").URL("role", role)
		if err != nil {
			t.Error(err)
		}

		url := server.URL + path.String()
		res, err := http.Post(url, "application/json", nil)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != http.StatusCreated {
			t.Errorf("Success expected at %s: %d", url, res.StatusCode)
		}
	}
}
