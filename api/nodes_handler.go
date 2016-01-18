package api

import "net/http"

func ListNodes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}
