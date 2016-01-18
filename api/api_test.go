package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var (
	server *httptest.Server
	router *mux.Router
	reader io.Reader
)

func init() {
	router = NewRouter()
	server = httptest.NewServer(router)
}

func TestListNodes(t *testing.T) {
	url := server.URL + "/api/nodes/"
	res, err := http.Get(url)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Success expected at %s: %d", url, res.StatusCode)
	}
}
