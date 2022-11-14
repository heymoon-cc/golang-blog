package controller

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var tags *map[string]string
var pages *map[string]string
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

func SetPages(newPages *map[string]string) {
	pages = newPages
}

func TagExists(tagCode string) bool {
	if _, ok := (*tags)[tagCode]; ok {
		return true
	}
	return false
}

func GetTagName(tagCode string) string {
	return (*tags)[tagCode]
}

func PageExists(pageName string) bool {
	if _, ok := (*pages)[pageName]; ok {
		return true
	}
	return false
}

func GetPageName(pageName string) string {
	return (*pages)[pageName]
}

type View struct {
	Header struct {
		Tags  *map[string]string
		Pages *map[string]string
		Title string
		Host  string
	}
	Page string
	Main interface{}
	Year int
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
	view.Header.Pages = pages
	view.Main = data
	view.Page = name
	view.Year = time.Now().Year()
	err = ts.Execute(w, view)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
