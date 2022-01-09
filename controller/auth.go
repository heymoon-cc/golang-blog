package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var username = os.Getenv("USER")
var password = os.Getenv("PASS")

func requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Restricted\"")
	w.WriteHeader(401)
}

func handleAuth(w http.ResponseWriter, r *http.Request, needRequestAuth bool) error {
	u, p, ok := r.BasicAuth()
	if !ok {
		if needRequestAuth {
			requestAuth(w)
		}
		return errors.New("auth required")
	}
	if u != username || p != password {
		fmt.Println(u)
		fmt.Println(p)
		fmt.Println(username)
		fmt.Println(password)
		if needRequestAuth {
			requestAuth(w)
		}
		return errors.New("unauthorized")
	}
	return nil
}
