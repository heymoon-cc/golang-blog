package controller

import (
  "github.com/gorilla/mux"
  "io"
  "mime"
  "net/http"
  "os"
  "path/filepath"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
  file := mux.Vars(r)["file"]
  data, err := os.Open("./static/dist/" + file)
  if err != nil {
    http.NotFound(w, r)
    return
  }
  w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
  _, err = io.Copy(w, data)
  if err != nil {
    http.NotFound(w, r)
    return
  }
}
