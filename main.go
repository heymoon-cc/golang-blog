package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"main/controller"
	"main/model"
	"net/http"
	"os"
)

type Config struct {
	Host  string            `json:"host"`
	Title string            `json:"title"`
	Tags  map[string]string `json:"tags"`
	Pages map[string]string `json:"pages"`
}

func main() {
	var config Config
	tagsFile, _ := os.ReadFile("./config.json")
	err := json.Unmarshal(tagsFile, &config)
	if err != nil {
		panic(err)
	}
	model.Connect()
	controller.SetHost(config.Host)
	controller.SetTitle(config.Title)
	controller.SetTags(&config.Tags)
	controller.SetPages(&config.Pages)
	server := &http.Server{
		Addr:         os.Getenv("ADDR"),
		Handler:      getRouter(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(server.ListenAndServe())
}
