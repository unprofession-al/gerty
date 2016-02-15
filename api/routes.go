package api

import "net/http"

type Route struct {
	// method
	m string
	// pattern
	p string
	// HandlerFunc
	h http.HandlerFunc
}

type Routes map[string]Route
