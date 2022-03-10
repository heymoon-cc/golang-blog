package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"html/template"
	"main/model"
	"net/http"
	"time"
)

type ArticleView struct {
	ID         string
	Title      string
	Content    template.HTML
	CreatedAt  time.Time
	Draft      bool
	Authorized bool
}

func createArticleView(article *model.Article, authorized bool) ArticleView {
	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})
	content := markdown.ToHTML([]byte(article.Content), nil, renderer)
	return ArticleView{
		article.ID.String(),
		article.Title,
		template.HTML(content),
		article.CreatedAt.Time(),
		article.Draft,
		authorized}
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	article := model.FindArticle(id)
	err = handleAuth(w, r, false)
	renderTemplate("./ui/html/article.page.tmpl", w, createArticleView(article, err == nil), "article")
}

func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	err := handleAuth(w, r, true)
	if err != nil {
		return
	}
	if r.Method == "GET" {
		renderTemplate("./ui/html/create.page.tmpl", w, nil, "create")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var article model.Article
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&article)
	if err != nil {
		return
	}
	model.CreateArticle(&article)
	http.Redirect(w, r, fmt.Sprintf("/article/%s", article.ID.String()), http.StatusFound)
}

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	err := handleAuth(w, r, true)
	if err != nil {
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	article := model.FindArticle(id)
	if r.Method == "GET" {
		renderTemplate("./ui/html/create.page.tmpl", w, article, "update")
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(article)
	if err != nil {
		fmt.Println(err)
		return
	}
	model.UpdateArticle(article)
	http.Redirect(w, r, fmt.Sprintf("/article/%s", article.ID.String()), http.StatusFound)
}
