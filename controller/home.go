package controller

import (
	"main/model"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	articles := model.AllArticles()
	var view TagView
	view.Tag = "All"
	for _, article := range *articles {
		view.Articles = append(view.Articles, createArticleView(&article, false))
	}
	renderTemplate("./ui/html/list.page.tmpl", w, &view, "index")
}
