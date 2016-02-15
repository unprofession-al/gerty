package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

type key int

const User key = 0

func setUserContext(req *http.Request, name string) {
	context.Set(req, User, name)
}

func getUserContext(req *http.Request) string {
	if user := context.Get(req, User); user != nil {
		return user.(string)
	}
	return ""
}

func userContextMiddleware(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 {
			setUserContext(req, "Sam Bell")
		} else {
			b, err := base64.StdEncoding.DecodeString(auth[1])
			if err != nil {
				panic("Crappy Auth Header Encoding")
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				panic("Crappy Auth Header")
			}

			setUserContext(req, pair[0])
		}
		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}
