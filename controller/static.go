package controller

import (
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func handleStatic(dir string, w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	data, err := os.Open(dir + "/" + file)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
	_, err = io.Copy(w, data)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	handleStatic("./static/dist", w, r)
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	handleStatic("./files", w, r)
}
