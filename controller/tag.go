package controller

import (
	"github.com/gorilla/mux"
	"main/model"
	"net/http"
)

type TagView struct {
	Tag      string
	Articles []ArticleView
}

func TagHandler(w http.ResponseWriter, r *http.Request) {
	tag := mux.Vars(r)["tag"]
	articles := model.ArticlesByTag(tag)
	var view TagView
	view.Tag = GetTagName(tag)
	for _, article := range *articles {
		view.Articles = append(view.Articles, createArticleView(&article, false))
	}
	renderTemplate("./ui/html/list.page.tmpl", w, &view, tag)
}
