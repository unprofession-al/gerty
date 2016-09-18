package api

import (
	"net/http"

	"github.com/unprofession-al/gerty/helpers"
	"github.com/unprofession-al/gerty/middleware"
)

func whoAmI(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(middleware.HeaderUserName)

	helpers.Respond(res, req, http.StatusOK, username)
}

func getNodeVarsProviders(res http.ResponseWriter, req *http.Request) {
	helpers.Respond(res, req, http.StatusOK, np)
}

func getConfig(res http.ResponseWriter, req *http.Request) {
	helpers.Respond(res, req, http.StatusOK, cfg)
}

func sitemapV1(res http.ResponseWriter, req *http.Request) {
	helpers.Respond(res, req, http.StatusOK, routes["v1"])
}
