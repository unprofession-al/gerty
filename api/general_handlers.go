package api

import "net/http"

func notImplemented(w http.ResponseWriter, r *http.Request) {
	user := GetUserContext(r)
	w.WriteHeader(http.StatusNotImplemented)
	out := "Function Not Yet Implemented, " + user + "\n"
	w.Write([]byte(out))
}
