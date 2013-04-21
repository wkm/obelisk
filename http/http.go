package main

import (
	_ "circuit/load"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Printf("obelisk/")
	http.HandleFunc("/zk/", zkHandler)
	http.HandleFunc("/host/", hostHandler)
	http.HandleFunc("/worker/", workerHandler)
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("err: %s", err)
	}
}

func respondError(rw http.ResponseWriter, msg string) {
	rw.Header().Set("Status-Code", "502")
	log.Printf("err: %s", msg)
	err := getTemplates().ExecuteTemplate(rw, "error.html", msg)
	if err != nil {
		fmt.Fprintf(rw, "Could not render error template: %s", err.Error())
	}
}

func hostHandler(rw http.ResponseWriter, req *http.Request) {
	renderTemplate(rw, "/host.html", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		path = "/index.html"
	}
	// try to find a template with the given name
	renderTemplate(rw, path, nil)
}
