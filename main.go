package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"main/controller"
	"main/model"
	"net/http"
	"os"
)

type Config struct {
	Host  string
	Title string
	Tags  map[string]string
}

func main() {
	var config Config
	tagsFile, _ := ioutil.ReadFile("./config.json")
	err := json.Unmarshal(tagsFile, &config)
	if err != nil {
		panic(err)
	}
	model.Connect()
	controller.SetHost(config.Host)
	controller.SetTitle(config.Title)
	controller.SetTags(&config.Tags)
	server := &http.Server{
		Addr:         os.Getenv("ADDR"),
		Handler:      getRouter(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(server.ListenAndServe())
}
