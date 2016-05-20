package middleware

import (
	"fmt"
	"net/http"
)

func WebHook(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(res, req)
		username := req.Header.Get(HeaderUserName)
		fmt.Println(username)
		res.Header().Set("X-Message", "Username: "+username)
	}

	return http.HandlerFunc(fn)
}
