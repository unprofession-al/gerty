package api

import "net/http"

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Function Not Yet Implemented\n"))
}
