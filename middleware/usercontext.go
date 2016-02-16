// Packege middleware contains all the required middlewares requried by gerty.
package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

const (
	HeaderUserName        = "G-UserName"
	DefaultUserName       = "Sam Bell"
	DefaultBrokenUserName = "Bam Sell"
)

// UserContext reads the 'Authorization' header from the request, decodes the
// credentials and stores the user name as new header 'G-UserName'
func UserContext(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 {
			req.Header.Set(HeaderUserName, DefaultUserName)
		} else {
			username, _ := extractUserName(auth[1])
			req.Header.Set(HeaderUserName, username)
		}
		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}

func extractUserName(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return DefaultBrokenUserName, errors.New("Credentials are not properly encoded")
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return DefaultBrokenUserName, errors.New("Decoded credentials are malformated")
	}

	return credentials[0], nil
}
