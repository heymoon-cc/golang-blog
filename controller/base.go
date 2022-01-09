package controller

import (
	"html/template"
	"log"
	"net/http"
)

var tags *map[string]string
var title string
var host string
var templates = make(map[string]*template.Template)

func SetTitle(newTitle string) {
	title = newTitle
}

func SetHost(newHost string) {
	host = newHost
}

func SetTags(newTags *map[string]string) {
	tags = newTags
}

func GetTagName(tagCode string) string {
	return (*tags)[tagCode]
}

type View struct {
	Header struct {
		Tags  *map[string]string
		Title string
		Host  string
	}
	Page string
	Main interface{}
}

func renderTemplate(url string, w http.ResponseWriter, data interface{}, name string) {
	var err error
	ts, found := templates[url]
	if !found {
		ts, err = template.ParseFiles(url, "./ui/html/base.layout.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		templates[url] = ts
	}
	view := View{}
	view.Header.Host = host
	view.Header.Title = title
	view.Header.Tags = tags
	view.Main = data
	view.Page = name
	err = ts.Execute(w, view)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
