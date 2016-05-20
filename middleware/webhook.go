package middleware

import "net/http"

func WebHook(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(res, req)
		username := req.Header.Get(HeaderUserName)
		res.Header().Set("X-Message", "Username: "+username)
	}

	return http.HandlerFunc(fn)
}
