package api

import "net/http"

func listNodes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}
