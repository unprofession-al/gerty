package middleware

import "net/http"

func CorsHeaders(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			res.Header().Set("Access-Control-Allow-Origin", origin)
		}
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, Authorization")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}
