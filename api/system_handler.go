package api

import (
	"net/http"

	"github.com/unprofession-al/gerty/helpers"
	"github.com/unprofession-al/gerty/middleware"
)

func whoAmI(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(middleware.HeaderUserName)

	helpers.Respond(res, req, http.StatusCreated, username)
}
