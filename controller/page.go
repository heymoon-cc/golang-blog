package controller

import (
	"github.com/gorilla/mux"
	"main/model"
	"net/http"
)

type PageView struct {
	Tag      string
	Articles []ArticleView
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	articles := model.AllArticlesByTag(page)
	var view TagView
	view.Tag = GetPageName(page)
	for _, article := range *articles {
		view.Articles = append(view.Articles, createArticleView(&article, false))
	}
	renderTemplate("./ui/html/page.page.tmpl", w, &view, page)
}
