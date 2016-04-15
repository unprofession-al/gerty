package middleware

import (
	"fmt"
	"net/http"
)

func CorsHeaders(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			res.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			res.Header().Set("Access-Control-Allow-Origin", "*")
		}
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, Authorization")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		if req.Method == http.MethodOptions {
			fmt.Fprintf(res, "Hello")
			return
		}
		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}
