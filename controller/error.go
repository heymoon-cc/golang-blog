package controller

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `/`, http.StatusTemporaryRedirect)
}
