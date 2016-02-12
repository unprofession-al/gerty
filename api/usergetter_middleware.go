package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

type key int

const User key = 0

type UserGetter struct {
}

func NewUserGetter() *UserGetter {
	return &UserGetter{}
}

func (l *UserGetter) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 {
		SetUserContext(req, "Sam Bell")
	} else {

		b, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			panic("Crappy Auth Header Encoding")
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			panic("Crappy Auth Header")
		}

		SetUserContext(req, pair[0])
	}
	next(res, req)
}

func SetUserContext(req *http.Request, name string) {
	context.Set(req, User, name)
}

func GetUserContext(req *http.Request) string {
	if user := context.Get(req, User); user != nil {
		return user.(string)
	}
	return ""
}
