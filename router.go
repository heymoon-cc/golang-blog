package main

import (
	"github.com/gorilla/mux"
	"main/controller"
)

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	r.HandleFunc("/tag/{tag}", controller.TagHandler)
	r.HandleFunc("/article/{id}", controller.ArticleHandler).Methods("GET")
	r.HandleFunc("/static/{file}", controller.StaticHandler).Methods("GET")
	r.HandleFunc("/files/{file}", controller.FilesHandler).Methods("GET")
	r.HandleFunc("/admin/create", controller.CreateArticleHandler).Methods("GET", "POST")
	r.HandleFunc("/admin/update/{id}", controller.UpdateArticleHandler).Methods("GET", "POST")
	r.HandleFunc("/{page}", controller.PageHandler)
	return r
}
